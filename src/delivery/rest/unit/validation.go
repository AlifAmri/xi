// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/delivery/model"
)

var (
	errInvalidPreparationUnit = "Dokumen tidak valid"
)

func validPreparationUnit(id int64) (r *model.PreparationUnit, e error) {
	e = orm.NewOrm().Raw("SELECT * FROM preparation_unit where id = ? and is_active = ?", id, 0).QueryRow(&r)

	return
}
