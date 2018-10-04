// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ptype

import (
	"git.qasico.com/cuxs/orm"
)

var (
	errRequiredName = "nama harus diisi"
	errUniqueName   = "nama tersebut telah digunakan"
	errInvalidID    = "type item tidak dapat ditemukan"
	errCascadeID    = "type masih digunakan oleh item"
)

func validName(name string, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM partnership_type where name = ? and id != ?", name, exclude).QueryRow(&total)

	return total == 0
}

func validID(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM partnership_type where id = ?", id).QueryRow(&total)

	return total == 1
}

func validDelete(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM partnership where type_id = ?", id).QueryRow(&total)

	return total == 0
}
