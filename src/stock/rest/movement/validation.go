// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package movement

import (
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	model2 "git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/stock/rest/unit"
	"git.qasico.com/gudang/api/src/warehouse"
)

var (
	errRequiredUnit         = "Stock Unit harus diisi"
	errRequiredDestination  = "Tujuan lokasi harus diisi"
	errRequiredQuantity     = "Quantity harus diisi"
	errRequiredContainer    = "Container harus diisi"
	errInvalidUnit          = "Stock unit tidak valid"
	errInvalidDestination   = "Tujuan lokasi tidak valid"
	errInvalidQuantity      = "Quantity melebihi stock yang ada"
	errInvalidStockMovement = "Movement tidak valid"
	errInvalidMergeUnit     = "Tidak dapat menyatukan stock dengan stock unit ini"
	errInvalidStockStorage  = "Storage tidak valid"
	errInvalidContainer     = "Container tidak valid"
	errInvalidOrigin        = "Lokasi asal tidak valid"
	errUniqueMovement       = "Stock unit sedang dalam proses movement"
	errLocationStockopname  = "Lokasi sedang dalam proses stockopname"
	errSameDestination      = "Tidak dapat melakukan movement pada lokasi yang sama"
)

func validUnit(ide string) (u *model.StockUnit, e error) {
	var id int64
	if id, e = common.Decrypt(ide); e == nil {
		u, e = unit.Show(id)
	}

	return
}

func validUniqueMovement(su *model.StockUnit) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(id) from stock_movement where (unit_id = ? or merge_unit = ?) and status != ?", su.ID, su.ID, "finish").QueryRow(&total)

	return total == 0
}

func validLocationMovement(l *warehouse.Location) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(id) from stock_opname where location_id = ? and status = ?", l.ID, "active").QueryRow(&total)

	return total == 0
}

func validMergeUnit(ide string, su *model.StockUnit, l *warehouse.Location) (u *model.StockUnit, e error) {
	var id int64
	if id, e = common.Decrypt(ide); e == nil {
		e = orm.NewOrm().Raw("SELECT * FROM stock_unit su "+
			"INNER JOIN stock_storage ss ON ss.id = su.storage_id "+
			"WHERE su.id = ? AND su.item_id = ? AND su.batch_id = ? and su.is_defect = ? "+
			"AND su.status = ? AND ss.location_id = ?;", id, su.Item.ID, su.Batch.ID, su.IsDefect, "stored", l.ID).QueryRow(&u)
	}

	return
}

func validLocation(ide string) (d *warehouse.Location, e error) {
	d = new(warehouse.Location)
	if d.ID, e = common.Decrypt(ide); e == nil {
		e = d.Read()
	}

	return
}

func validStorage(ide string) (d *model.StockStorage, e error) {
	d = new(model.StockStorage)
	if d.ID, e = common.Decrypt(ide); e == nil {
		e = d.Read()
	}

	return
}

func validContainer(ide string) (c *model2.Item, e error) {
	var ID int64

	if ID, e = common.Decrypt(ide); e == nil {
		e = orm.NewOrm().Raw("SELECT i.* from item i INNER JOIN item_type it on i.type_id = it.id where i.id = ? and i.is_active = ? and it.is_container = ?", ID, 1, 1).QueryRow(&c)
	}

	return
}

func validStockMovement(id int64, status string) (m *model.StockMovement, e error) {
	m = new(model.StockMovement)
	o := orm.NewOrm()
	e = o.QueryTable(m).Filter("id", id).Filter("status", status).RelatedSel().Limit(1).One(m)

	return
}
