// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"fmt"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/delivery/model"
)

type item struct {
	ID        string  `json:"id"`
	ItemCode  string  `json:"item_code" valid:"required"`
	BatchCode string  `json:"batch_code"`
	Quantity  float64 `json:"quantity" valid:"required|gte:1"`

	PreparationPlanItem *model.PreparationPlanItem `json:"-"`
}

func (rp *item) Validate(index int, o *validation.Output) {
	var e error

	if rp.BatchCode != "" {
		if rp.BatchCode, e = validBatchCode(rp.BatchCode); e != nil {
			o.Failure(fmt.Sprintf("items.%d.batch_code.invalid", index), errInvalidBatchCode)
		}
	}

	if rp.ID != "" {
		if rp.PreparationPlanItem, e = validPreparationPlanItem(rp.ID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.id.invalid", index), errInvalidReceivingPlan)
		}
	}
}

func (rp *item) Save(r *model.PreparationPlan) {
	rpp := &model.PreparationPlanItem{
		Plan:      r,
		ItemCode:  rp.ItemCode,
		BatchCode: rp.BatchCode,
		Quantity:  rp.Quantity,
	}

	if rp.PreparationPlanItem != nil {
		rpp.ID = rp.PreparationPlanItem.ID
	}

	rpp.Save()
}
