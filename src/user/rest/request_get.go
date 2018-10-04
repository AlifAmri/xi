// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/orm"
)

// Get get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]user.User, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(user.User))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []user.User
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// Show find a single data user using field and value condition.
func Show(id int64, permission bool) (u *user.User, e error) {
	u = new(user.User)
	o := orm.NewOrm()

	if e = o.QueryTable(u).Filter("id", id).RelatedSel().Limit(1).One(u); e == nil && permission {

		o.Raw("SELECT p.id as identity, p.*, (pu.id IS NOT NULL) AS grantted FROM privilege p "+
			"left join privilege_user pu on pu.privilege_id = p.id and pu.user_id = ? "+
			"where p.is_active = 1;", u.ID).QueryRows(&u.Privileges)
	}

	return
}
