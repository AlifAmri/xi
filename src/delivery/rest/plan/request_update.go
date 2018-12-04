// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"time"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID           int64   `json:"-" valid:"required"`
	PartnerID    string  `json:"partner_id"`
	DocumentCode string  `json:"document_code" valid:"required"`
	ShippedDate  string  `json:"shipped_date" valid:"required"`
	Note         string  `json:"note"`
	Items        []*item `json:"items" valid:"required"`

	Session         *auth.SessionData      `json:"-"`
	PreparationPlan *model.PreparationPlan `json:"-"`
	Partner         *model2.Partnership    `json:"-"`
	ShippedAt       time.Time              `json:"-"`
	TotalQuantity   float64                `json:"-"`
}

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.PreparationPlan, e = validPreparationPlan(ur.ID, "draft"); e != nil {
		o.Failure("id.invalid", errInvalidReceiptPlan)
	}

	if ur.PartnerID != "" {
		if ur.Partner, e = validPartner(ur.PartnerID); e != nil {
			o.Failure("partner_id.invalid", errInvalidPartner)
		}
	}

	if ur.DocumentCode != "" && ur.PreparationPlan != nil {
		if !validDocumentCode(ur.DocumentCode, ur.PreparationPlan.ID) {
			o.Failure("document_code.unique", errUniqueCode)
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
			ur.TotalQuantity += item.Quantity
		}
	}

	return o
}

func (ur *updateRequest) Messages() map[string]string {
	return map[string]string{
		"document_code.required": errRequiredDocumentCode,
		"shipped_date.required":  errRequiredReceivingAt,
		"items.required":         errRequiredItems,
	}
}

func (ur *updateRequest) Save() (u *model.PreparationPlan, e error) {
	u = &model.PreparationPlan{
		ID:            ur.PreparationPlan.ID,
		Status:        "draft",
		DocumentCode:  ur.DocumentCode,
		TotalQuantity: ur.TotalQuantity,
		Partner:       ur.Partner,
		Note:          ur.Note,
		CreatedBy:     ur.Session.User.(*user.User),
		CreatedAt:     ur.PreparationPlan.CreatedAt,
		ShippedAt:     ur.ShippedAt,
	}

	if e = u.Save(); e == nil {
		orm.NewOrm().LoadRelated(ur.PreparationPlan, "Items", 0)

		for _, item := range ur.Items {
			item.Save(ur.PreparationPlan)
		}

		// delete item yang tidak lagi di proses
		if ur.PreparationPlan.Items != nil {
			for _, plan := range ur.PreparationPlan.Items {
				var used bool
				for _, item := range ur.Items {
					if item.PreparationPlanItem != nil && plan.ID == item.PreparationPlanItem.ID {
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
