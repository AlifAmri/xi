// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

var (
	errRequiredReceiving    = "Receiving dokumen harus diisi"
	errRequiredItem         = "Item harus diisi"
	errRequiredBatchCode    = "Kode Batch harus diisi"
	errRequiredQuantity     = "Quantity harus diisi"
	errRequiredCheckedBy    = "User harus diisi"
	errRequiredUnit         = "Unit harus diisi"
	errRequiredLocation     = "Lokasi harus diisi"
	errInvalidReceiving     = "Receiving dokumen tidak valid"
	errInvalidCheckedBy     = "User tidak valid"
	errInvalidUnit          = "Unit tidak valid atau code telah digunakan"
	errInvalidItem          = "Item tidak valid"
	errInvalidReceivingUnit = "Dokumen tidak valid"
	errInvalidLocation      = "Lokasi tidak valid"
)

func validReceiving(ide string) (r *model.Receiving, e error) {
	r = new(model.Receiving)
	if r.ID, e = common.Decrypt(ide); e == nil {
		e = r.Read()
	}

	return
}

func validCheckedBy(ide string) (r *user.User, e error) {
	r = new(user.User)
	if r.ID, e = common.Decrypt(ide); e == nil {
		e = r.Read()
	}

	return
}

func validUnitCode(code string, exclude int64) bool {
	var total int64
	o := orm.NewOrm()
	o.Raw("SELECT count(*) FROM stock_unit where code = ?", code).QueryRow(&total)
	if total == 0 {
		o.Raw("SELECT count(*) FROM receiving_unit where unit_code = ? and id != ?", code, exclude).QueryRow(&total)
	}

	return total == 0
}

func validLocation(ide string) (r *warehouse.Location, e error) {
	r = new(warehouse.Location)
	if r.ID, e = common.Decrypt(ide); e == nil {
		e = r.Read()
	}

	return
}

func validReceivingUnit(id int64, s string) (r *model.ReceivingUnit, e error) {
	e = orm.NewOrm().Raw("SELECT * FROM receiving_unit where id = ? and is_active = ?", id, 0).QueryRow(&r)

	return
}
