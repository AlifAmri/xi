// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/vehicle/model"
)

// Show find a single data incoming_vehicles using field and value condition.
func Show(id int64) (*model.IncomingVehicle, error) {
	m := new(model.IncomingVehicle)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Get get all data incoming_vehicles that matched with query request parameters.
// returning slices of IncomingVehicles, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]model.IncomingVehicle, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.IncomingVehicle))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.IncomingVehicle
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
