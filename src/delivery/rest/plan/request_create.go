// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"time"

	"git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/partnership/model"

	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	PartnerID    string  `json:"partner_id"`
	DocumentCode string  `json:"document_code" valid:"required"`
	ShippedDate  string  `json:"shipped_date" valid:"required"`
	Note         string  `json:"note"`
	Items        []*item `json:"items" valid:"required"`

	Session       *auth.SessionData   `json:"-"`
	Partner       *model2.Partnership `json:"-"`
	ShippedAt     time.Time           `json:"-"`
	TotalQuantity float64             `json:"-"`
}

func (cr *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.DocumentCode != "" {
		if !validDocumentCode(cr.DocumentCode, 0) {
			o.Failure("document_code.unique", errUniqueCode)
		}
	}

	if cr.PartnerID != "" {
		if cr.Partner, e = validPartner(cr.PartnerID); e != nil {
			o.Failure("partner_id.invalid", errInvalidPartner)
		}
	}

	if cr.ShippedDate != "" {
		if cr.ShippedAt, e = time.Parse("2006-01-02", cr.ShippedDate); e != nil {
			o.Failure("shipped_date.invalid", errInvalidDate)
		}
	}

	if len(cr.Items) > 0 {
		for i, item := range cr.Items {
			item.Validate(i, o)
			cr.TotalQuantity += item.Quantity
		}
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
	return map[string]string{
		"document_code.required": errRequiredDocumentCode,
		"shipped_date.required":  errRequiredReceivingAt,
		"items.required":         errRequiredItems,
	}
}

func (cr *createRequest) Save() (u *model.PreparationPlan, e error) {
	u = &model.PreparationPlan{
		Status:        "draft",
		DocumentCode:  cr.DocumentCode,
		TotalQuantity: cr.TotalQuantity,
		Partner:       cr.Partner,
		Note:          cr.Note,
		CreatedBy:     cr.Session.User.(*user.User),
		CreatedAt:     time.Now(),
		ShippedAt:     cr.ShippedAt,
	}

	if e = u.Save(); e == nil {
		for _, item := range cr.Items {
			item.Save(u)
		}
	}

	return
}
