// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/usergroup"

	"git.qasico.com/cuxs/orm"
)

// Show find a single data usergroup using field and value condition.
func Show(id int64) (ug *usergroup.Usergroup, err error) {
	ug = new(usergroup.Usergroup)
	o := orm.NewOrm()

	if err = o.QueryTable(ug).Filter("id", id).Limit(1).One(ug); err == nil {

		o.Raw("SELECT p.id as identity, p.*, (pu.id IS NOT NULL) AS grantted FROM privilege p "+
			"left join privilege_usergroup pu on pu.privilege_id = p.id and pu.usergroup_id = ? "+
			"where p.is_active = 1;", ug.ID).QueryRows(&ug.Privileges)
	}

	return
}

// Get get all data usergroup that matched with query request parameters.
// returning slices of Usergroup, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]usergroup.Usergroup, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(usergroup.Usergroup))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []usergroup.Usergroup
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
