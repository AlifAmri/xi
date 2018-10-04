// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package location

import (
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/orm"
)

// Show find a single data warehouse_Location using field and value condition.
func Show(id int64) (*warehouse.Location, error) {
	m := new(warehouse.Location)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Get get all data warehouse_Location that matched with query request parameters.
// returning slices of WarehouseLocation, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]warehouse.Location, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(warehouse.Location))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []warehouse.Location
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
