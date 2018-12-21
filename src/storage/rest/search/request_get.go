// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package search

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/warehouse"
)

func findLocation(itemCode string, batchCode string, year string) (m []*warehouse.Location, total int64, err error) {
	var d []*warehouse.Location
	o := orm.NewOrm()
	if batchCode != "" {
		total, err = o.Raw("SELECT wl.* FROM stock_unit su "+
			"inner join stock_storage ss on ss.id = su.storage_id "+
			"inner join warehouse_location wl on wl.id = ss.location_id "+
			"inner join item i on i.id = su.item_id "+
			"inner join item_batch ib on ib.id = su.batch_id "+
			"where i.code = ? and ib.code = ? group by location_id;", itemCode, batchCode).QueryRows(&d)
	}else if batchCode == "" && year != ""{
		total, err = o.Raw("SELECT wl.* FROM stock_unit su "+
			"inner join stock_storage ss on ss.id = su.storage_id "+
			"inner join warehouse_location wl on wl.id = ss.location_id "+
			"inner join item i on i.id = su.item_id "+
			"inner join item_batch ib on ib.id = su.batch_id "+
			"where i.code = ? and ib.code LIKE ? group by location_id;", itemCode, "%"+year).QueryRows(&d)
	}else {
		total, err = o.Raw("SELECT wl.* FROM stock_unit su "+
			"inner join stock_storage ss on ss.id = su.storage_id "+
			"inner join warehouse_location wl on wl.id = ss.location_id "+
			"inner join item i on i.id = su.item_id "+
			"inner join item_batch ib on ib.id = su.batch_id "+
			"where i.code = ? group by location_id;", itemCode).QueryRows(&d)
	}

	return d, total, err
}
