// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uom

import (
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
)

var (
	errRequiredName = "nama harus diisi"
	errRequiredType = "tipe item harus dipilih"
	errUniqueName   = "nama tersebut telah digunakan"
	errInvalidID    = "uom item tidak dapat ditemukan"
	errInvalidType  = "tipe item tidak dapat ditemukan"
	errCascadeID    = "uom masih digunakan oleh item"
)

func validName(name string, typeID int64, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item_uom where name = ? and type_id = ? and id != ?", name, typeID, exclude).QueryRow(&total)

	return total == 0
}

func validID(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item_uom where id = ?", id).QueryRow(&total)

	return total == 1
}

func validDelete(id int64) bool {
	// @todo validasi untuk memperbolehkan menghapus uom
	return id != -1
}

func validType(ide string) (it *model.ItemType, e error) {
	it = new(model.ItemType)
	if it.ID, e = common.Decrypt(ide); e == nil {
		e = it.Read()
	}

	return
}
