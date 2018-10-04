// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package group

import (
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
)

var (
	errRequiredName = "nama harus diisi"
	errRequiredType = "tipe item harus dipilih"
	errUniqueName   = "nama tersebut telah digunakan"
	errInvalidID    = "kategory item tidak dapat ditemukan"
	errInvalidType  = "tipe item tidak dapat ditemukan"
	errCascadeID    = "kategory masih digunakan oleh item"
)

type attribute struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

func validName(name string, typeID int64, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item_group where name = ? and id != ? and type_id = ?", name, exclude, typeID).QueryRow(&total)

	return total == 0
}

func validID(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item_group where id = ?", id).QueryRow(&total)

	return total == 1
}

func validDelete(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item where group_id = ?", id).QueryRow(&total)

	return total == 0
}

func validType(ide string) (it *model.ItemType, e error) {
	it = new(model.ItemType)
	if it.ID, e = common.Decrypt(ide); e == nil {
		e = it.Read()
	}

	return
}
