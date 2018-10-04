// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/storage/model"
)

// Show find a single data usergroup using field and value condition.
func Show(id int64) (ug *model.StorageGroup, err error) {
	ug = new(model.StorageGroup)
	o := orm.NewOrm()

	if err = o.QueryTable(ug).Filter("id", id).Limit(1).One(ug); err == nil {

		if ug.Type == "item_category" {
			o.Raw("select * from item_category where name = ?", ug.TypeValue).QueryRow(&ug.ItemCategory)
		} else if ug.Type == "item_group" {
			o.Raw("select * from item_group where name = ?", ug.TypeValue).QueryRow(&ug.ItemGroup)
		}

		o.Raw("SELECT p.id as identity, p.*, (pu.id IS NOT NULL) AS is_selected FROM warehouse_area p "+
			"left join storage_group_area pu on pu.warehouse_area_id = p.id and pu.storage_group_id = ? "+
			"where p.is_active = 1;", ug.ID).QueryRows(&ug.Areas)
	}

	return
}

// Get get all data usergroup that matched with query request parameters.
// returning slices of Usergroup, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]*model.StorageGroup, total int64, err error) {
	// make new orm query
	q, o := rq.Query(new(model.StorageGroup))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.StorageGroup
	if _, err = q.All(&mx, rq.Fields...); err == nil {

		for _, ig := range mx {
			o.Raw("SELECT p.id as identity, p.*, (pu.id IS NOT NULL) AS is_selected FROM warehouse_area p "+
				"inner join storage_group_area pu on pu.warehouse_area_id = p.id and pu.storage_group_id = ? "+
				"where p.is_active = 1;", ig.ID).QueryRows(&ig.Areas)
		}

		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
