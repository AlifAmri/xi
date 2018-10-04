// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package item

import (
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/warehouse"
)

var (
	errRequiredCode     = "kode harus diisi"
	errRequiredName     = "nama harus diisi"
	errRequiredType     = "tipe item harus dipilih"
	errUniqueCode       = "kode tersebut telah digunakan"
	errInvalidGroup     = "grup item tidak dapat ditemukan"
	errInvalidType      = "tipe item tidak dapat ditemukan"
	errInvalidCategory  = "kategori item tidak dapat ditemukan"
	errInvalidArea      = "area tidak dapat ditemukan"
	errInvalidBarcode   = "barcode number salah"
	errInvalidID        = "item tidak dapat ditemukan"
	errCascadeID        = "kategory masih digunakan oleh item"
	errAlreadyDeactived = "status item sudah tidak aktif"
	errAlreadyActived   = "status item sudah aktif"
)

type attribute struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

func validCode(code string, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item where code = ? and id != ?", code, exclude).QueryRow(&total)

	return total == 0
}

func validItem(id int64) (ig *model.Item, e error) {
	ig = &model.Item{ID: id}
	e = ig.Read()

	return
}

func validGroup(ide string) (ig *model.ItemGroup, e error) {
	ig = new(model.ItemGroup)
	if ig.ID, e = common.Decrypt(ide); e == nil {
		e = ig.Read()
	}

	return
}

func validType(ide string) (it *model.ItemType, e error) {
	it = new(model.ItemType)
	if it.ID, e = common.Decrypt(ide); e == nil {
		e = it.Read()
	}

	return
}

func validCategory(ide string) (ic *model.ItemCategory, e error) {
	ic = new(model.ItemCategory)
	if ic.ID, e = common.Decrypt(ide); e == nil {
		e = ic.Read()
	}

	return
}

func validArea(ide string) (wa *warehouse.Area, e error) {
	wa = new(warehouse.Area)
	if wa.ID, e = common.Decrypt(ide); e == nil {
		e = wa.Read()
	}

	return
}

func validID(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item where id = ?", id).QueryRow(&total)

	return total == 1
}

func validDelete(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM item_batch where item_id = ?", id).QueryRow(&total)

	return total == 0
}
