// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"errors"
	"fmt"
	"strconv"

	"git.qasico.com/gudang/api/src/receiving/model"
	model2 "git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	model3 "git.qasico.com/gudang/api/src/inventory/model"
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
	errInvalidBatchCode     = "Format kode batch tidak valid"
	errInvalidItemCode      = "Kode tidak ditemukan atau tidak aktif"
	errLocationFull 		= "Lokasi sudah Full"
)

func validBatchCode(code string) (c string, e error) {
	c = code

	if len(c) == 3 {
		c = fmt.Sprintf("%s%s", "0", code)
	}

	if len(c) != 4 {
		e = errors.New("wrong format")
	} else {
		cx := c[0:2]
		if !validWeek(cx) {
			return "", errors.New("not valid")
		}
	}

	return
}

func validWeek(s string) bool {
	i, e := strconv.Atoi(s)
	if e == nil && i > 0 && i < 54 {
		return true
	}

	return false
}

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

func validItemCode(code string) bool {
	var total int64

	orm.NewOrm().Raw("SELECT count(*) FROM item where code = ? and type_id = ? and is_active = ?", code, 1, 1).QueryRow(&total)

	return total == 1
}

func countLocationMoved(id int64) (total int) {
	orm.NewOrm().Raw("SELECT count(distinct ss.id) FROM stock_storage ss" +
		" inner join stock_unit su on su.storage_id = ss.id" +
		" where ss.location_id = ? and su.status = ?", id,"stored").QueryRow(&total)
	return
}

func countMovement(id int64) (total int) {
	orm.NewOrm().Raw("SELECT count(id) FROM stock_movement where destination_id = ? and is_merger = ? and is_not_full = ? and status != ?", id, 0, 0,"finish").QueryRow(&total)
	return
}

func validUnitCode(code string, exclude int64, rp *model.Receiving) (r *model2.StockUnit, e error) {
	var total int64
	o := orm.NewOrm()

	o.Raw("SELECT count(*) FROM receiving_unit where unit_code = ? and id != ?", code, exclude).QueryRow(&total)
	if total == 0 {
		// cek apakah ada di receiving plan
		if rp != nil {
			if rp.Plan != nil {
				o.Raw("SELECT count(*) FROM receipt_plan_item where unit_code = ? and plan_id = ?", code, rp.Plan.ID).QueryRow(&total)
			}

			if total == 0 {
				o.Raw("SELECT count(*) FROM receiving_document rd "+
					"inner join stock_unit su on su.id = rd.unit_id "+
					"where rd.receiving_id = ? and su.code = ?;", rp.ID, code).QueryRow(&total)
			}

			if total > 0 {
				o.Raw("SELECT * FROM stock_unit where code = ?", code).QueryRow(&r)
			}
		}

		if r == nil {
			o.Raw("SELECT count(*) FROM stock_unit where code = ?", code).QueryRow(&total)
			if total > 0 {
				e = errors.New("already used")
			}
		}
	} else {
		e = errors.New("already receipt")
	}

	return
}

func validLocation(ide string) (r *warehouse.Location, e error) {
	r = new(warehouse.Location)
	if r.ID, e = common.Decrypt(ide); e == nil {
		e = r.Read()
	}

	return
}

func validPallet(ide string) (r *model3.Item, e error) {
	var ID int64
	if ID, e = common.Decrypt(ide); e == nil {
		e = orm.NewOrm().Raw("SELECT * FROM item i where i.id = ? and i.is_active = ? and i.type_id = ?", ID, 1, 3).QueryRow(&r)
	}

	return
}
func validReceivingUnit(id int64, s string) (r *model.ReceivingUnit, e error) {
	//e = orm.NewOrm().Raw("SELECT * FROM receiving_unit where id = ? and is_active = ?", id, 0).QueryRow(&r)
	r = new(model.ReceivingUnit)
	o := orm.NewOrm()
	if e = o.QueryTable(r).Filter("id", id).Filter("is_active", 0).RelatedSel().Limit(1).One(r); e != nil {
		return nil, e
	}

	return
}
