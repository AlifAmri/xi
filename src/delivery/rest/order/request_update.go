// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/user"
)

type updateRequest struct {
	ID              int64  `json:"-" valid:"required"`
	PartnerID       string `json:"partner_id"`
	NumberContainer string `json:"number_container"`
	NumberSeal      string `json:"number_seal"`
	Note            string `json:"note"`
	CheckerID       string `json:"checker_id"`
	SupplyID        string `json:"supply_id"`
	DocLoc          string `json:"docloc"`

	Items []*item `json:"items" valid:"required"`

	Session       *auth.SessionData    `json:"-"`
	DeliveryOrder *model.DeliveryOrder `json:"-"`
	Partner       *model2.Partnership  `json:"-"`
	Supply        *user.User           `json:"-"`
	Chekcer       *user.User           `json:"-"`
}

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.DeliveryOrder, e = validDeliveryOrder(ur.ID, "active"); e != nil {
		o.Failure("id.invalid", errInvalidDeliveryOrder)
	}

	if ur.PartnerID != "" {
		if ur.Partner, e = validPartner(ur.PartnerID); e != nil {
			o.Failure("partner_id.invalid", errInvalidPartner)
		}
	}

	if ur.SupplyID != "" {
		if ur.Supply, e = validUser(ur.SupplyID); e != nil {
			o.Failure("supply_id.invalid", errInvalidPartner)
		}
	}
	if ur.CheckerID != "" {
		if ur.Chekcer, e = validUser(ur.CheckerID); e != nil {
			o.Failure("checker_id.invalid", errInvalidPartner)
		}
	}
	if len(ur.Items) > 0 && ur.DeliveryOrder != nil {
		for i, item := range ur.Items {
			item.Validate(i, o, ur.DeliveryOrder)
		}
	}

	return o
}

func (ur *updateRequest) Messages() map[string]string {
	return map[string]string{
		"items.required": errRequiredItems,
	}
}

func (ur *updateRequest) Save() (u *model.DeliveryOrder, e error) {
	u = ur.DeliveryOrder

	u.Partner = ur.Partner
	u.NumberContainer = ur.NumberContainer
	u.NumberSeal = ur.NumberSeal
	u.Note = ur.Note
	u.Supply = ur.Supply
	u.Checker = ur.Chekcer
	u.DocLoc = ur.DocLoc

	if e = u.Save("partner_id", "number_container", "number_seal", "note", "supply_id", "checker_id", "docloc"); e == nil {
		orm.NewOrm().LoadRelated(ur.DeliveryOrder, "Items", 0)

		for _, item := range ur.Items {
			item.Save(ur.DeliveryOrder)
		}

		// delete item yang tidak lagi di proses
		if ur.DeliveryOrder.Items != nil {
			for _, plan := range ur.DeliveryOrder.Items {
				var used bool
				for _, item := range ur.Items {
					if item.Preparation != nil && plan.ID == item.Preparation.ID {
						used = true
					}
				}

				if !used {
					plan.DeliveryOrder = nil
					plan.Save("delivery_order_id")
				}
			}
			// update incoming vehicle status when DO is filled
			go func() {
				o := orm.NewOrm()
				o.Raw("UPDATE incoming_vehicle ivh INNER JOIN delivery_order do ON do.vehicle_id = ivh.id "+
					"SET ivh.status = 'in_progress' WHERE do.id = ? AND ivh.status = ?", ur.ID, "in_queue").Exec()
			}()
		}

	}

	return
}
