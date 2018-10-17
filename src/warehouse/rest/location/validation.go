// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package location

import (
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

var (
	errRequiredCode     = "Kode lokasi harus diisi"
	errRequiredName     = "Nama lokasi harus diisi"
	errRequiredArea     = "Area harus dipilih"
	errRequiredCordX    = "Titik X harus diisi"
	errRequiredCordY    = "Titik Y harus diisi"
	errRequiredCordW    = "Lebar Titik harus diisi"
	errRequiredCordH    = "Tinggi Titik harus diisi"
	errRequiredCapacity = "Kapasitas lokasi harus diisi"
	errUniqueCode       = "Kode lokasi tersebut telah digunakan"
	errInvalidLocation  = "Lokasi tidak dapat ditemukan"
	errInvalidArea      = "Area tidak dapat ditemukan"
	errAlreadyActived   = "Status area sudah aktif"
	errAlreadyDeactived = "Status area sudah tidak aktif"
)

func validCode(code string, areaID int64, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM warehouse_location where code = ? and warehouse_area_id = ? and id != ?", code, areaID, exclude).QueryRow(&total)

	return total == 0
}

func validArea(ide string) (wa *warehouse.Area, e error) {
	wa = new(warehouse.Area)
	if wa.ID, e = common.Decrypt(ide); e == nil {
		e = wa.Read()
	}

	return
}

func validLocation(id int64) (u *warehouse.Location, e error) {
	u = &warehouse.Location{ID: id}
	e = u.Read()

	return
}
