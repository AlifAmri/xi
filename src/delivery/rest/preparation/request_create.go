// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/user"
	"time"

	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	PartnerID    string  `json:"partner_id" valid:"required"`
	DocumentCode string  `json:"document_code" valid:"required"`
	DocumentFile string  `json:"document_file"`
	ShippedDate  string  `json:"shipped_date" valid:"required"`
	Note         string  `json:"note"`
	Items        []*item `json:"items" valid:"required"`

	Session           *auth.SessionData   `json:"-"`
	Partner           *model2.Partnership `json:"-"`
	ShippedAt         time.Time           `json:"-"`
	TotalQuantityPlan float64             `json:"-"`
}

func (ur *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.DocumentCode != "" {
		if !validDocumentCode(ur.DocumentCode, 0) {
			o.Failure("document_code.unique", errUniqueCode)
		}
	}

	if ur.PartnerID != "" {
		if ur.Partner, e = validPartner(ur.PartnerID); e != nil {
			o.Failure("partner_id.invalid", errInvalidPartner)
		}
	}

	if ur.ShippedDate != "" {
		if ur.ShippedAt, e = time.Parse("2006-01-02", ur.ShippedDate); e != nil {
			o.Failure("shipped_date.invalid", errInvalidDate)
		}
	}

	if len(ur.Items) > 0 {
		for i, item := range ur.Items {
			item.Validate(i, o)
			ur.TotalQuantityPlan += item.Quantity
		}
	}

	return o
}

func (ur *createRequest) Messages() map[string]string {
	return map[string]string{
		"partner_id.required":    errRequiredPartner,
		"document_code.required": errRequiredDocumentCode,
		"items.required":         errRequiredItems,
	}
}

func (ur *createRequest) Save() (u *model.Preparation, e error) {
	u = &model.Preparation{
		Partner:           ur.Partner,
		Supervisor:        ur.Session.User.(*user.User),
		Status:            "draft",
		DocumentCode:      ur.DocumentCode,
		DocumentFile:      ur.DocumentFile,
		TotalQuantityPlan: ur.TotalQuantityPlan,
		Note:              ur.Note,
		CreatedBy:         ur.Session.User.(*user.User),
		CreatedAt:         time.Now(),
		ShippedAt:         ur.ShippedAt,
	}

	if e = u.Save(); e == nil {
		for _, item := range ur.Items {
			item.Save(u)
		}
	}

	return
}
