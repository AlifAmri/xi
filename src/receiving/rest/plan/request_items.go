// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"fmt"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/receiving/model"
)

type item struct {
	ID        string  `json:"id"`
	ItemCode  string  `json:"item_code" valid:"required"`
	BatchCode string  `json:"batch_code" valid:"required"`
	UnitCode  string  `json:"unit_code"`
	Quantity  float64 `json:"quantity" valid:"required,min:1"`

	ReceiptPlanItem *model.ReceiptPlanItem `json:"-"`
}

func (rp *item) Validate(index int, o *validation.Output) {
	var e error

	if rp.ID != "" {
		if rp.ReceiptPlanItem, e = validReceiptPlanItem(rp.ID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.id.invalid", index), errInvalidReceivingPlan)
		}
	}
}

func (rp *item) Save(r *model.ReceiptPlan) {
	rpp := &model.ReceiptPlanItem{
		Plan:      r,
		UnitCode:  rp.UnitCode,
		ItemCode:  rp.ItemCode,
		BatchCode: rp.BatchCode,
		Quantity:  rp.Quantity,
	}

	if rp.ReceiptPlanItem != nil {
		rpp.ID = rp.ReceiptPlanItem.ID
	}

	rpp.Save()
}
