// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package opname

import (
	"errors"
	"fmt"
	"strconv"

	iModel "git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

var (
	errRequiredType           = "Type dokumen harus dipilih"
	errRequiredItems          = "Item harus diisi"
	errRequiredItem           = "Item harus diisi"
	errRequiredLocation       = "Lokasi harus diisi"
	errRequiredBatch          = "Kode batch harus diisi"
	errRequiredQuantity       = "Quantity harus diisi atau gunakan void"
	errRequiredContainerNum   = "No. container harus diisi"
	errInvalidStockOpname     = "Dokumen tidak valid"
	errInvalidLocation        = "Lokasi tidak valid atau stockopname masih aktif untuk lokasi ini"
	errInvalidItem            = "Item tidak valid"
	errInvalidUnit            = "Kode unit tidak valid"
	errInvalidStockopnameItem = "Item tidak valid"
	errInvalidContainer       = "Container tidak valid"
	errInvalidBatchCode       = "Format kode batch tidak valid"
	errAlreadyFinished        = "Dokumen telah selesai"
)

func validItem(ide string) (i *iModel.Item, e error) {
	i = new(iModel.Item)
	if i.ID, e = common.Decrypt(ide); e == nil {
		if e = i.Read(); e == nil {
			i.Type.Read()
		}

		if i.IsActive == int8(0) {
			e = errors.New("item not active")
		}
	}

	return
}

func validUnit(ide string) (u *model.StockUnit, e error) {
	u = new(model.StockUnit)
	if u.ID, e = common.Decrypt(ide); e == nil {
		e = u.Read()
	}

	return
}

func validUnitCode(code string) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM stock_unit where code = ?", code).QueryRow(&total)

	return total == 0
}

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

func validLocation(ide string) (l *warehouse.Location, e error) {
	l = new(warehouse.Location)

	if l.ID, e = common.Decrypt(ide); e == nil {
		if e = l.Read(); e == nil {
			if l.IsActive == 0 {
				e = errors.New("location is not active")
			} else {
				// cek stockopname lain dengan location ini
				var total int64
				orm.NewOrm().Raw("SELECT count(id) from stock_opname where location_id = ? and status = ?", l.ID, "active").QueryRow(&total)
				if total > 0 {
					e = errors.New("location still have stockopname")
				}
			}
		}
	}

	return
}

func validStockOpname(id int64) (so *model.StockOpname, e error) {
	e = orm.NewOrm().Raw("SELECT * from stock_opname where id = ? and status = ?", id, "active").QueryRow(&so)

	return
}

func validStockopnameItem(ide string) (soi *model.StockOpnameItem, e error) {
	soi = new(model.StockOpnameItem)
	if soi.ID, e = common.Decrypt(ide); e == nil {
		e = soi.Read()
	}

	return
}

func validContainer(ide string) (c *iModel.Item, e error) {
	var ID int64

	if ID, e = common.Decrypt(ide); e == nil {
		e = orm.NewOrm().Raw("SELECT i.* from item i INNER JOIN item_type it on i.type_id = it.id where i.id = ? and i.is_active = ? and it.is_container = ?", ID, 1, 1).QueryRow(&c)
	}

	return
}
