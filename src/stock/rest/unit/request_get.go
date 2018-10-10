// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/cuxs/orm"
	model3 "git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/stock/model"
)

// Show find a single data item_type using field and value condition.
func Show(id int64) (*model.StockUnit, error) {
	m := new(model.StockUnit)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Get get all data item_type that matched with query request parameters.
// returning slices of ItemType, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]model.StockUnit, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.StockUnit))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.StockUnit
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func historyMovement(id int64, rq *orm.RequestQuery) (m *[]model.StockMovement, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.StockMovement))

	nc := orm.NewCondition()
	cond2 := q.GetCond().AndCond(nc.And("unit_id", id).Or("new_unit", id).Or("merge_unit", id))
	q = q.SetCond(cond2)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.StockMovement
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func historyStockopname(id int64, rq *orm.RequestQuery) (m *[]model.StockOpnameItem, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.StockOpnameItem))

	q = q.Filter("unit_id", id)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.StockOpnameItem
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func historyReceiving(id int64, rq *orm.RequestQuery) (m *[]model2.ReceivingUnit, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model2.ReceivingUnit))

	q = q.Filter("unit_id", id)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model2.ReceivingUnit
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func historyPreparation(id int64, rq *orm.RequestQuery) (m *[]model3.PreparationUnit, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model3.PreparationUnit))

	q = q.Filter("unit_id", id)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model3.PreparationUnit
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
