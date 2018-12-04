// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"fmt"
	"time"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/auth"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID           int64   `json:"-" valid:"required"`
	PartnerID    string  `json:"partner_id"`
	DocumentCode string  `json:"document_code" valid:"required"`
	ReceivedDate string  `json:"received_date" valid:"required"`
	Note         string  `json:"note"`
	Items        []*item `json:"items" valid:"required"`

	Session       *auth.SessionData   `json:"-"`
	ReceiptPlan   *model.ReceiptPlan  `json:"-"`
	Partner       *model2.Partnership `json:"-"`
	ReceivedAt    time.Time           `json:"-"`
	TotalQuantity float64             `json:"-"`
}

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.ReceiptPlan, e = validReceiptPlan(ur.ID, "draft"); e != nil {
		o.Failure("id.invalid", errInvalidReceiptPlan)
	}

	if ur.PartnerID != "" {
		if ur.Partner, e = validPartner(ur.PartnerID); e != nil {
			o.Failure("partner_id.invalid", errInvalidPartner)
		}
	}

	if ur.DocumentCode != "" && ur.ReceiptPlan != nil {
		if !validDocumentCode(ur.DocumentCode, ur.ReceiptPlan.ID) {
			o.Failure("document_code.unique", errUniqueCode)
		}
	}

	if ur.ReceivedDate != "" {
		if ur.ReceivedAt, e = time.Parse("2006-01-02", ur.ReceivedDate); e != nil {
			o.Failure("receiving_at.invalid", errInvalidDate)
		}
	}

	if len(ur.Items) > 0 {
		unitCode := make(map[string]int)
		for i, item := range ur.Items {
			item.Validate(i, o)
			ur.TotalQuantity += item.Quantity

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

func (ur *updateRequest) Messages() map[string]string {
	return map[string]string{
		"document_code.required": errRequiredDocumentCode,
		"receiving_at.required":  errRequiredReceivingAt,
		"items.required":         errRequiredItems,
	}
}

func (ur *updateRequest) Save() (u *model.ReceiptPlan, e error) {
	u = &model.ReceiptPlan{
		ID:            ur.ReceiptPlan.ID,
		Status:        "draft",
		DocumentCode:  ur.DocumentCode,
		TotalQuantity: ur.TotalQuantity,
		Partner:       ur.Partner,
		Note:          ur.Note,
		CreatedBy:     ur.Session.User.(*user.User),
		CreatedAt:     ur.ReceiptPlan.CreatedAt,
		ReceivedAt:    ur.ReceivedAt,
	}

	if e = u.Save(); e == nil {
		orm.NewOrm().LoadRelated(ur.ReceiptPlan, "Items", 0)

		for _, item := range ur.Items {
			item.Save(ur.ReceiptPlan)
		}

		// delete item yang tidak lagi di proses
		if ur.ReceiptPlan.Items != nil {
			for _, plan := range ur.ReceiptPlan.Items {
				var used bool
				for _, item := range ur.Items {
					if item.ReceiptPlanItem != nil && plan.ID == item.ReceiptPlanItem.ID {
						used = true
					}
				}

				if !used {
					plan.Delete()
				}
			}
		}
	}

	return u, e
}
