// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package item

import (
	"fmt"
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
)

type equationRequest struct {
	MasterID string      `json:"master_id" valid:"required"`
	Items    []*equation `json:"items" valid:"required"`

	Item    *model.Item       `json:"-"`
	Session *auth.SessionData `json:"-"`
}

type equation struct {
	ItemID string      `json:"item_id"`
	Item   *model.Item `json:"-"`
}

func (eq *equation) Validate(index int, o *validation.Output) {
	var e error

	if eq.ItemID != "" {
		if eq.Item, e = validItemID(eq.ItemID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.item_id.invalid", index), errInvalidItem)
		}
	} else {
		o.Failure(fmt.Sprintf("items.%d.item_id.required", index), errRequiredItem)
	}
}

func (ar *equationRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ar.MasterID != "" {
		if ar.Item, e = validItemID(ar.MasterID); e != nil {
			o.Failure("master_id.invalid", errInvalidID)
		}
	}

	if len(ar.Items) > 0 {
		for i, item := range ar.Items {
			item.Validate(i, o)
		}
	}

	return o
}

func (ar *equationRequest) Messages() map[string]string {
	return map[string]string{
		"master_id.required": errRequiredItem,
		"items.required":     errRequiredItem,
	}
}

func (ar *equationRequest) Save() error {
	// update item dengan equation sebelumnya
	if ar.Item.Equation != "" {
		orm.NewOrm().Raw("UPDATE item set equation = null where equation = ?", ar.Item.Equation).Exec()
	}

	ar.Item.Equation = common.RandomStr(10)

	for _, item := range ar.Items {
		item.Item.Equation = ar.Item.Equation
		item.Item.Save("equation")
	}

	return ar.Item.Save("equation")
}
