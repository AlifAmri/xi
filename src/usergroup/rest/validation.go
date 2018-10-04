// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import "git.qasico.com/cuxs/orm"

func validName(name string, exlude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM usergroup where name = ? and id != ?", name, exlude).QueryRow(&total)

	return total == 0
}

func validID(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM usergroup where id = ?", id).QueryRow(&total)

	return total == 1
}

func validDelete(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM user where usergroup_id = ?", id).QueryRow(&total)

	return total == 0
}
