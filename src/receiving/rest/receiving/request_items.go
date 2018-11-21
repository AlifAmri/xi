// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"fmt"

	ModelInventory "git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/receiving/model"
	ModelStock "git.qasico.com/gudang/api/src/stock/model"

	"git.qasico.com/cuxs/validation"
)

type item struct {
	ID        string  `json:"id"`
	ItemCode  string  `json:"item_code" valid:"required"`
	BatchCode string  `json:"batch_code"`
	UnitCode  string  `json:"unit_code"`
	Quantity  float64 `json:"quantity" valid:"required|gte:1"`

	Item              *ModelInventory.Item      `json:"-"`
	ItemBatch         *ModelInventory.ItemBatch `json:"-"`
	StockUnit         *ModelStock.StockUnit     `json:"-"`
	ReceivingDocument *model.ReceivingDocument  `json:"-"`
	IsNewItem         int8                      `json:"-"`
}

func (rp *item) Validate(index int, o *validation.Output) {
	var e error

	if rp.ItemCode != "" {
		if rp.Item, rp.IsNewItem, e = validItemCode(rp.ItemCode); e != nil {
			o.Failure(fmt.Sprintf("items.%d.item_code.invalid", index), errInvalidItemCode)
		}
	}

	if rp.Item != nil && rp.BatchCode != "" {
		if rp.BatchCode, e = validBatchCodeString(rp.BatchCode); e != nil {
			o.Failure(fmt.Sprintf("items.%d.batch_code.invalid", index), errInvalidBatchCode)
		} else {
			if rp.ItemBatch, e = validBatchCode(rp.BatchCode, rp.Item); e != nil {
				o.Failure(fmt.Sprintf("items.%d.batch_code.invalid", index), errInvalidBatchCode)
			}
		}
	}

	if rp.ItemBatch != nil && rp.UnitCode != "" {
		if rp.StockUnit, e = validStockUnit(rp.UnitCode, rp.Item, rp.ItemBatch); e != nil {
			o.Failure(fmt.Sprintf("items.%d.unit_code.invalid", index), errInvalidUnitCode)
		}
	}

	if rp.ItemBatch == nil && rp.UnitCode != "" {
		o.Failure(fmt.Sprintf("items.%d.batch_code.invalid", index), errInvalidBatchCode)
	}

	if rp.ID != "" {
		if rp.ReceivingDocument, e = validReceivingPlan(rp.ID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.id.invalid", index), errInvalidReceivingPlan)
		}
	}
}

func (rp *item) Save(r *model.Receiving) *model.ReceivingDocument {
	rpp := &model.ReceivingDocument{
		Receiving: r,
		Item:      rp.Item,
		Batch:     rp.ItemBatch,
		Unit:      rp.StockUnit,
		IsNew:     rp.IsNewItem,
		Quantity:  rp.Quantity,
	}

	if rp.ReceivingDocument != nil {
		rpp.ID = rp.ReceivingDocument.ID
		rpp.IsNew = rp.ReceivingDocument.IsNew
	}

	rpp.Save()

	return rpp
}
