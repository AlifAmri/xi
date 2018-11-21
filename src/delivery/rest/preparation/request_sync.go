// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"time"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	"git.qasico.com/gudang/api/src/user"
)

type syncRequest struct {
	PlanID string `json:"plan_id" valid:"required"`

	PreparationPlan   *model.PreparationPlan `json:"-"`
	Session           *auth.SessionData      `json:"-"`
	Items             []*item                `json:"-"`
	TotalQuantityPlan float64                `json:"-"`
}

func (cr *syncRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.PlanID != "" {
		if cr.PreparationPlan, e = validPreparationPlan(cr.PlanID); e != nil {
			o.Failure("plan_id.invalid", errInvalidPreparationPlan)
		}
	}

	if cr.PreparationPlan != nil {
		orm.NewOrm().LoadRelated(cr.PreparationPlan, "Items", 0)
		if cr.PreparationPlan.Items != nil {
			for n, i := range cr.PreparationPlan.Items {
				i := &item{
					ItemCode:  i.ItemCode,
					BatchCode: i.BatchCode,
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
		"plan_id.required": errRequiredPreparationPlan,
	}
}

func (cr *syncRequest) Save() (e error) {
	d := &model.Preparation{
		DocumentCode:      cr.PreparationPlan.DocumentCode,
		Partner:           cr.PreparationPlan.Partner,
		Plan:              cr.PreparationPlan,
		TotalQuantityPlan: cr.TotalQuantityPlan,
		CreatedBy:         cr.Session.User.(*user.User),
		CreatedAt:         time.Now(),
		ShippedAt:         cr.PreparationPlan.ShippedAt,
		Status:            "draft",
	}

	if e = d.Save(); e == nil {
		for _, item := range cr.Items {
			item.Save(d)
		}

		cr.PreparationPlan.Status = "active"
		cr.PreparationPlan.Save("status")
	}

	return
}
