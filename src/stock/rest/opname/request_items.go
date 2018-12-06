// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package opname

import (
	"fmt"
	"time"

	"git.qasico.com/gudang/api/src/inventory"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/validation"
	iModel "git.qasico.com/gudang/api/src/inventory/model"
)

type opnameItem struct {
	ID           string  `json:"id"`
	UnitID       string  `json:"unit_id"`
	ItemID       string  `json:"item_id" valid:"required"`
	ContainerID  string  `json:"container_id"`
	ContainerNum int8    `json:"container_num"`
	BatchCode    string  `json:"batch_code"`
	UnitCode     string  `json:"unit_code"`
	Quantity     float64 `json:"quantity"`
	Note         string  `json:"note"`
	IsDefect     int8    `json:"is_defect"`
	IsVoid       int8    `json:"is_void"`

	Item            *iModel.Item           `json:"-"`
	Container       *iModel.Item           `json:"-"`
	ItemBatch       *iModel.ItemBatch      `json:"-"`
	StockUnit       *model.StockUnit       `json:"-"`
	StockopnameItem *model.StockOpnameItem `json:"-"`
	IsNewUnit       int8                   `json:"-"`
}

func (oi *opnameItem) Validate(index int, o *validation.Output) {
	var e error

	// cek item dan batch apabila
	// unit id tidak diisi, karena stockopname
	// bisa membuat unitid baru berdasarkan item dan batch
	if oi.UnitID == "" {
		if oi.ItemID != "" {
			if oi.Item, e = validItem(oi.ItemID); e != nil {
				o.Failure(fmt.Sprintf("items.%d.item_id.invalid", index), errInvalidItem)
			}
		} else {
			o.Failure(fmt.Sprintf("items.%d.item_id.required", index), errRequiredItem)
		}

		if oi.BatchCode != "" && oi.Item != nil && oi.Item.ID > 0 {
			if oi.BatchCode, e = validBatchCode(oi.BatchCode); e == nil {
				oi.ItemBatch = inventory.GetBatch(oi.Item.ID, oi.BatchCode)
			} else {
				o.Failure(fmt.Sprintf("items.%d.batch_code.invalid", index), errInvalidBatchCode)
			}
		}

		if oi.UnitCode != "" && !validUnitCode(oi.UnitCode) {
			o.Failure(fmt.Sprintf("items.%d.unit_id.invalid", index), errInvalidUnit)
		}

		// kalau type item is_batch nya true, maka stockopname
		// item ini memerlukan batch
		if oi.Item != nil && oi.Item.ID > 0 && oi.Item.Type.IsBatch == int8(1) && oi.ItemBatch == nil {
			o.Failure(fmt.Sprintf("items.%d.batch_code.required", index), errRequiredBatch)
		}
	} else {
		if oi.StockUnit, e = validUnit(oi.UnitID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.unit_id.invalid", index), errInvalidUnit)
		}
	}

	if oi.Quantity == 0 {
		o.Failure(fmt.Sprintf("items.%d.quantity.invalid", index), errRequiredQuantity)
	}

	if oi.ContainerID != "" {
		if oi.Container, e = validContainer(oi.ContainerID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.container_id.invalid", index), errInvalidContainer)
		}

		if oi.ContainerNum == 0 {
			o.Failure(fmt.Sprintf("items.%d.container_num.invalid", index), errRequiredContainerNum)
		}
	}

	if oi.ID != "" {
		if oi.StockopnameItem, e = validStockopnameItem(oi.ID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.id.invalid", index), errInvalidStockopnameItem)
		}
	}
}

func (oi *opnameItem) Save(so *model.StockOpname, u *user.User) {

	if oi.StockUnit == nil && oi.Item.Type.IsBatch == int8(1) {
		oi.StockUnit = &model.StockUnit{
			Code:       oi.UnitCode,
			Item:       oi.Item,
			Batch:      oi.ItemBatch,
			Status:     "draft",
			IsDefect:   oi.IsDefect,
			ReceivedAt: time.Now(),
			CreatedBy:  u,
		}

		oi.StockUnit.GenerateCode("")
		err := oi.StockUnit.Save()
		oi.IsNewUnit = 1
		fmt.Println("show error unit -------------- :", err)
	}

	if oi.IsVoid == 1 {
		oi.Quantity = 0
		oi.Container = nil
		oi.ContainerNum = 0
		oi.IsDefect = 0
		oi.IsNewUnit = 0
	}

	soi := &model.StockOpnameItem{
		StockOpname:    so,
		Item:           oi.StockUnit.Item,
		Unit:           oi.StockUnit,
		UnitQuantity:   oi.StockUnit.Stock,
		ActualQuantity: oi.Quantity,
		Container:      oi.Container,
		ContrainerNum:  oi.ContainerNum,
		IsNewUnit:      oi.IsNewUnit,
		IsDefect:       oi.IsDefect,
		IsVoid:         oi.IsVoid,
		Note:           oi.Note,
	}

	if oi.StockopnameItem != nil {
		soi.ID = oi.StockopnameItem.ID
		soi.IsNewUnit = oi.StockopnameItem.IsNewUnit
	}

	er := soi.Save()
	fmt.Println("show error op item-------------- :", er)
}
