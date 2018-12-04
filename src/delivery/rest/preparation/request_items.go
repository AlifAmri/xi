// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"fmt"

	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/delivery/model"
	ModelInventory "git.qasico.com/gudang/api/src/inventory/model"
)

type item struct {
	ID        string  `json:"id"`
	ItemCode  string  `json:"item_code" valid:"required"`
	BatchCode string  `json:"batch_code"`
	Quantity  float64 `json:"quantity" valid:"required|gte:1"`

	Item                *ModelInventory.Item       `json:"-"`
	ItemBatch           *ModelInventory.ItemBatch  `json:"-"`
	PreparationDocument *model.PreparationDocument `json:"-"`
	IsNewItem           int8                       `json:"-"`
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

	if rp.ID != "" {
		if rp.PreparationDocument, e = validPreparationDocument(rp.ID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.id.invalid", index), errInvalidPreparationDocument)
		}
	}
}

func (rp *item) Save(r *model.Preparation) *model.PreparationDocument {
	rpp := &model.PreparationDocument{
		Preparation: r,
		Item:        rp.Item,
		Batch:       rp.ItemBatch,
		Quantity:    rp.Quantity,
	}

	if rp.PreparationDocument != nil {
		rpp.ID = rp.PreparationDocument.ID
	}

	rpp.Save()

	return rpp
}
