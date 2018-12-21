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
	errRequiredCode        = "Kode lokasi harus diisi"
	errRequiredName        = "Nama lokasi harus diisi"
	errRequiredArea        = "Area harus dipilih"
	errRequiredCordX       = "Titik X harus diisi"
	errRequiredCordY       = "Titik Y harus diisi"
	errRequiredCordW       = "Lebar Titik harus diisi"
	errRequiredCordH       = "Tinggi Titik harus diisi"
	errRequiredCapacity    = "Kapasitas lokasi harus diisi"
	errUniqueCode          = "Kode lokasi tersebut telah digunakan"
	errInvalidLocation     = "Lokasi tidak dapat ditemukan"
	errInvalidArea         = "Area tidak dapat ditemukan"
	errAlreadyActived      = "Status area sudah aktif"
	errAlreadyDeactived    = "Status area sudah tidak aktif"
	errCascadeID           = "Lokasi tidak dapat dihapus"
	errInvalidLocationArea = "Area di lokasi ini tidak dapat diganti"
)

func validLocationEmpty(locationID int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM stock_unit su "+
		"INNER JOIN stock_storage ss ON ss.id = su.storage_id "+
		"INNER JOIN warehouse_location wl ON wl.id = ss.location_id "+
		"WHERE wl.id = ? AND su.status = ? ", locationID, "stored").QueryRow(&total)

	return total == int64(0)
}

func validMovementLocation(locationID int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM stock_movement sm "+
		"INNER JOIN warehouse_location wl ON wl.id = sm.destination_id "+
		"WHERE wl.id = ? AND sm.status != ? ", locationID, "finish").QueryRow(&total)

	return total == int64(0)
}

func validStockopnameLocation(locationID int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM stock_opname so "+
		"INNER JOIN warehouse_location wl ON wl.id = so.location_id "+
		"WHERE wl.id = ? AND so.status = ? ", locationID, "active").QueryRow(&total)

	return total == int64(0)
}

func validPreparationLocation(locationID int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM preparation p "+
		"INNER JOIN warehouse_location wl ON wl.id = p.location_id "+
		"WHERE wl.id = ? AND p.status != ? ", locationID, "finish").QueryRow(&total)

	return total == int64(0)
}

func validReceivingLocation(locationID int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM receiving_unit ru "+
		"INNER JOIN warehouse_location wl ON wl.id = ru.location_received "+
		"INNER JOIN receiving r ON r.id = ru.receiving_id "+
		"WHERE wl.id = ? AND r.status = ? ", locationID, "active").QueryRow(&total)

	return total == int64(0)
}

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

func validDelete(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM stock_storage where location_id = ?", id).QueryRow(&total)

	return total == 0
}
