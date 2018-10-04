// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/receiving/services"
)

type syncRequest struct {
	ID     int64  `json:"-" valid:"required"`
	PlanID string `json:"plan_id" valid:"required"`

	Receiving         *model.Receiving   `json:"-"`
	ReceiptPlan       *model.ReceiptPlan `json:"-"`
	Session           *auth.SessionData  `json:"-"`
	Items             []*item            `json:"-"`
	TotalQuantityPlan float64            `json:"-"`
}

func (cr *syncRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.Receiving, e = validSyncID(cr.ID); e != nil {
		o.Failure("id.invalid", errInvalidReceiving)
	}

	if cr.PlanID != "" {
		if cr.ReceiptPlan, e = validReceiptPlan(cr.PlanID); e != nil {
			o.Failure("plan_id.invalid", errInvalidReceiptPlan)
		}
	}

	if cr.ReceiptPlan != nil {
		orm.NewOrm().LoadRelated(cr.ReceiptPlan, "Items", 0)
		if cr.ReceiptPlan.Items != nil {
			for n, i := range cr.ReceiptPlan.Items {
				i := &item{
					ItemCode:  i.ItemCode,
					BatchCode: i.BatchCode,
					UnitCode:  i.UnitCode,
					Quantity:  i.Quantity,
				}

				i.Validate(n, o)

				cr.TotalQuantityPlan += i.Quantity
				cr.Items = append(cr.Items, i)
			}
		}
	}

	return o
}

func (cr *syncRequest) Messages() map[string]string {
	return map[string]string{
		"plan_id.required": errRequiredReceiptPlan,
	}
}

func (cr *syncRequest) Save() (e error) {
	cr.Receiving.DocumentCode = cr.ReceiptPlan.DocumentCode
	cr.Receiving.Partner = cr.ReceiptPlan.Partner
	cr.Receiving.Plan = cr.ReceiptPlan
	cr.Receiving.TotalQuantityPlan = cr.TotalQuantityPlan

	if e = cr.Receiving.Save("plan_id", "document_code", "partner_id", "total_quantity_plan"); e == nil {
		for _, item := range cr.Items {
			item.Save(cr.Receiving)
		}

		cr.ReceiptPlan.Status = "active"
		cr.ReceiptPlan.Save("status")

		go services.CreateActual(cr.Receiving)
	}

	return
}
