// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package area

import (
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/orm"
)

// Show find a single data warehouse_area using field and value condition.
func Show(id int64) (*warehouse.Area, error) {
	m := new(warehouse.Area)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Get get all data warehouse_area that matched with query request parameters.
// returning slices of Area, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]warehouse.Area, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(warehouse.Area))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []warehouse.Area
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
