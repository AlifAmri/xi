// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"fmt"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"time"

	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	PartnerID    string  `json:"partner_id"`
	DocumentCode string  `json:"document_code" valid:"required"`
	ReceivedDate string  `json:"received_date" valid:"required"`
	Note         string  `json:"note"`
	Items        []*item `json:"items" valid:"required"`

	Session       *auth.SessionData   `json:"-"`
	Partner       *model2.Partnership `json:"-"`
	ReceivedAt    time.Time           `json:"-"`
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

	if cr.ReceivedDate != "" {
		if cr.ReceivedAt, e = time.Parse("2006-01-02", cr.ReceivedDate); e != nil {
			o.Failure("received_date.invalid", errInvalidDate)
		}
	}

	if len(cr.Items) > 0 {
		// unit code kalau diisi harus unique
		unitCode := make(map[string]int)
		for i, item := range cr.Items {
			item.Validate(i, o)
			cr.TotalQuantity += item.Quantity

			if item.UnitCode != "" {
				if _, ok := unitCode[item.UnitCode]; ok {
					o.Failure(fmt.Sprintf("items.%d.unit_code.invalid", i), errUniqueUnit)
				} else {
					unitCode[item.UnitCode] = 1
				}
			}
		}
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
	return map[string]string{
		"document_code.required": errRequiredDocumentCode,
		"received_date.required": errRequiredReceivingAt,
		"items.required":         errRequiredItems,
	}
}

func (cr *createRequest) Save() (u *model.ReceiptPlan, e error) {
	u = &model.ReceiptPlan{
		Status:        "draft",
		DocumentCode:  cr.DocumentCode,
		TotalQuantity: cr.TotalQuantity,
		Partner:       cr.Partner,
		Note:          cr.Note,
		CreatedBy:     cr.Session.User.(*user.User),
		CreatedAt:     time.Now(),
		ReceivedAt:    cr.ReceivedAt,
	}

	if e = u.Save(); e == nil {
		for _, item := range cr.Items {
			item.Save(u)
		}
	}

	return
}
