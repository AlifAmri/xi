// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uom

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
)

// Show find a single data item_uom using field and value condition.
func Show(id int64) (*model.ItemUom, error) {
	m := new(model.ItemUom)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Get get all data item_uom that matched with query request parameters.
// returning slices of ItemUom, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]model.ItemUom, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.ItemUom))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.ItemUom
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
