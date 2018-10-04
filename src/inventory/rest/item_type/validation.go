// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package itype

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
)

var (
	errRequiredName     = "nama harus diisi"
	errUniqueName       = "nama tersebut telah digunakan"
	errInvalidID        = "type item tidak dapat ditemukan"
	errInvalidContainer = "hanya satu type item yang dapat menjadi container"
	errCascadeID        = "type masih digunakan oleh item"
	errStaticBatch      = "tidak dapat menrubah batch"
	errStaticContainer  = "tidak dapat menrubah container"
)

func validName(name string, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item_type where name = ? and id != ?", name, exclude).QueryRow(&total)

	return total == 0
}

func validID(id int64) (it *model.ItemType, e error) {
	it = new(model.ItemType)
	it.ID = id

	e = it.Read()

	return
}

func validContainer(exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item_type where is_container = ? and id != ?", 1, exclude).QueryRow(&total)

	return total == 0
}

func validDelete(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item where type_id = ?", id).QueryRow(&total)
	if total == 0 {
		orm.NewOrm().Raw("SELECT count(*) FROM item_group where type_id = ?", id).QueryRow(&total)
		if total == 0 {
			orm.NewOrm().Raw("SELECT count(*) FROM item_category where type_id = ?", id).QueryRow(&total)
			if total == 0 {
				orm.NewOrm().Raw("SELECT count(*) FROM item_uom where type_id = ?", id).QueryRow(&total)
			}
		}
	}

	return total == 0
}
