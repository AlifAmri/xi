// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/stock/model"
)

// GetLogItem get all data warehouse_area that matched with query request parameters.
// returning slices of Area, total data without limit and error.
func GetLogItem(rq *orm.RequestQuery, id int64) (m *[]model.StockLog, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.StockLog))

	q = q.Filter("item_id", id)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.StockLog
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetLogBatch get all data warehouse_area that matched with query request parameters.
// returning slices of Area, total data without limit and error.
func GetLogBatch(rq *orm.RequestQuery, id int64) (m *[]model.StockLog, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.StockLog))

	q = q.Filter("batch_id", id)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.StockLog
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetLogUnit get all data warehouse_area that matched with query request parameters.
// returning slices of Area, total data without limit and error.
func GetLogUnit(rq *orm.RequestQuery, id int64) (m *[]model.StockLog, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.StockLog))

	q = q.Filter("stock_unit_id", id)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.StockLog
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
