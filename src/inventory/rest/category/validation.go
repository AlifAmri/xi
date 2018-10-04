// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package category

import (
	"errors"
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
)

var (
	errRequiredName      = "nama harus diisi"
	errRequiredType      = "tipe item harus dipilih"
	errUniqueName        = "nama tersebut telah digunakan"
	errInvalidID         = "kategory item tidak dapat ditemukan"
	errInvalidType       = "tipe item tidak dapat ditemukan"
	errInvalidParent     = "parent kategori tidak ditemukan"
	errInvalidSelfParent = "tidak dapat memilih kategori ini"
	errCascadeID         = "kategory masih digunakan oleh item"
)

func validName(name string, typeID int64, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item_category where name = ? and type_id = ? and id != ?", name, typeID, exclude).QueryRow(&total)

	return total == 0
}

func validID(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item_category where id = ?", id).QueryRow(&total)

	return total == 1
}

func validDelete(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item where category_id = ?", id).QueryRow(&total)

	return total == 0
}

func validType(ide string) (it *model.ItemType, e error) {
	it = new(model.ItemType)
	if it.ID, e = common.Decrypt(ide); e == nil {
		e = it.Read()
	}

	return
}

func validParent(ide string, typeID int64) (p *model.ItemCategory, e error) {
	p = new(model.ItemCategory)
	if p.ID, e = common.Decrypt(ide); e == nil {
		e = p.Read()

		if p.Type.ID != typeID {
			e = errors.New("different item type")
		}
	}

	return
}
