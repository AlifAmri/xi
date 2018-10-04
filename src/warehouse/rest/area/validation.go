// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package area

import (
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/orm"
)

var (
	errRequiredCode     = "Kode area harus diisi"
	errRequiredName     = "Nama area harus diisi"
	errRequiredType     = "Tipe area harus diisi"
	errUniqueCode       = "Kode area tersebut telah digunakan"
	errInvalidArea      = "Area tidak dapat ditemukan"
	errInvalidType      = "Tipe area tidak valid"
	errAlreadyActived   = "Status area sudah aktif"
	errAlreadyDeactived = "Status area sudah tidak aktif"
)

func validCode(code string, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM warehouse_area where code = ? and id != ?", code, exclude).QueryRow(&total)

	return total == 0
}

func validArea(id int64) (u *warehouse.Area, e error) {
	u = &warehouse.Area{ID: id}
	e = u.Read()

	return
}
