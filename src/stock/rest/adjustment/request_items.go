// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package adjustment

import (
	"fmt"
	"git.qasico.com/gudang/api/src/stock/model"

	"git.qasico.com/cuxs/validation"
	iModel "git.qasico.com/gudang/api/src/inventory/model"
)

type opnameItem struct {
	ID       string  `json:"id"`
	ItemID   string  `json:"item_id" valid:"required"`
	Quantity float64 `json:"quantity" valid:"required"`
	Note     string  `json:"note"`
	IsDefect int8    `json:"is_defect"`

	Item            *iModel.Item           `json:"-"`
	StockopnameItem *model.StockOpnameItem `json:"-"`
}

func (oi *opnameItem) Validate(index int, o *validation.Output) {
	var e error

	if oi.ItemID != "" {
		if oi.Item, e = validItem(oi.ItemID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.item_id.invalid", index), errInvalidItem)
		}
	} else {
		o.Failure(fmt.Sprintf("items.%d.item_id.required", index), errRequiredItem)
	}

	if oi.Quantity == 0 {
		o.Failure(fmt.Sprintf("items.%d.quantity.required", index), errRequiredQuantity)
	}

	if oi.Item != nil && oi.Quantity < 0 {
		if oi.Item.Stock < (oi.Quantity * -1) {
			o.Failure(fmt.Sprintf("items.%d.quantity.required", index), errInvalidQuantity)
		}
	}

	if oi.ID != "" {
		if oi.StockopnameItem, e = validStockopnameItem(oi.ID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.item_id.invalid", index), errInvalidStockopnameItem)
		}
	}
}

func (oi *opnameItem) Save(so *model.StockOpname) {
	soi := &model.StockOpnameItem{
		StockOpname:    so,
		Item:           oi.Item,
		UnitQuantity:   oi.Item.Stock,
		ActualQuantity: oi.Item.Stock + oi.Quantity,
		Note:           oi.Note,
	}

	if oi.StockopnameItem != nil {
		soi.ID = oi.StockopnameItem.ID
	}

	e := soi.Save()
	fmt.Println(e)
	fmt.Println("=============")
}
