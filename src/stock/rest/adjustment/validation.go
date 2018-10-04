// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package adjustment

import (
	"errors"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	iModel "git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/stock/model"
)

var (
	errRequiredItems          = "Item harus diisi"
	errRequiredItem           = "Item harus diisi"
	errRequiredQuantity       = "Quantity tidak boleh 0"
	errInvalidStockOpname     = "Dokumen tidak valid"
	errInvalidLocation        = "Lokasi tidak valid"
	errInvalidItem            = "Item tidak valid"
	errInvalidStockopnameItem = "Item tidak valid"
	errAlreadyFinished        = "Dokumen telah selesai"
	errInvalidQuantity        = "Quantity tidak mencukupi"
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

func validStockOpname(id int64) (so *model.StockOpname, e error) {
	e = orm.NewOrm().Raw("SELECT * from stock_opname where id = ? and status = ? and type = ?", id, "active", "adjustment").QueryRow(&so)

	return
}

func validStockopnameItem(ide string) (soi *model.StockOpnameItem, e error) {
	soi = new(model.StockOpnameItem)
	if soi.ID, e = common.Decrypt(ide); e == nil {
		e = soi.Read()
	}

	return
}
