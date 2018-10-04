// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/privilege"

	"git.qasico.com/cuxs/orm"
)

// Get all data privilege that matched with query request parameters.
// returning slices of Privilege, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]privilege.Privilege, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(privilege.Privilege))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []privilege.Privilege
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
