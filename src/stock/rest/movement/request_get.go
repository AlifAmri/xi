// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package movement

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/stock/model"
)

// Show find a single data warehouse_area using field and value condition.
func Show(id int64) (*model.StockMovement, error) {
	m := new(model.StockMovement)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "Items", 2)
	return m, nil
}

// Get get all data warehouse_area that matched with query request parameters.
// returning slices of Area, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]model.StockMovement, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.StockMovement))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.StockMovement
	if _, err = q.RelatedSel("pallet_id").All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
