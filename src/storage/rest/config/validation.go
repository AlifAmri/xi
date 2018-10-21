// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"errors"
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/storage/model"
	"git.qasico.com/gudang/api/src/warehouse"
)

var (
	errRequiredName        = "Nama harus diisi"
	errRequiredType        = "Tipe harus dipilih"
	errRequiredValue       = "Value harus diisi"
	errUniqueName          = "Nama tersebut telah digunakan"
	errInvalidStorageGroup = "Grup storage tidak dapat ditemukan"
	errInvalidArea         = "Area tidak valid"
	errInvalidLocation     = "Lokasi tidak valid"
	errAlreadyActived      = "Grup storage sudah aktif"
	errAlreadyDeactived    = "Grup storage sudah tidak aktif"
	errRequiredItem        = "Area dan lokasi harus diisi"
	errUniqueLocation      = "Lokasi sudah dipakai untuk konfig lain"
)

func validName(name string, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM storage_group where name = ? and id != ?", name, exclude).QueryRow(&total)

	return total == 0
}

func validID(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM storage_group where id = ?", id).QueryRow(&total)

	return total == 1
}

func validStorageGroup(id int64) (ug *model.StorageGroup, e error) {
	ug = new(model.StorageGroup)
	ug.ID = id

	e = ug.Read()

	return
}

func validArea(ide string) (u *warehouse.Area, e error) {
	u = new(warehouse.Area)
	if u.ID, e = common.Decrypt(ide); e == nil {
		u.Type = "storage"
		e = u.Read("id", "type")
	}

	return
}

func validLocationFrom(ide string, area *warehouse.Area) (u *warehouse.Location, e error) {
	u = new(warehouse.Location)
	if u.ID, e = common.Decrypt(ide); e == nil {
		u.Area = area
		u.IsActive = 1
		e = u.Read("id", "warehouse_area_id", "is_active")
	}

	return
}

func validLocationEnd(ide string, area *warehouse.Area, from *warehouse.Location) (u *warehouse.Location, e error) {
	u = new(warehouse.Location)
	if u.ID, e = common.Decrypt(ide); e == nil {
		u.Area = area
		u.IsActive = 1
		if e = u.Read("id", "warehouse_area_id", "is_active"); e == nil {
			if u.ID < from.ID {
				e = errors.New("end lokasi lebih kecil dari end")
			}
		}
	}

	return
}

func validLocation(wl *warehouse.Location, exID int64) bool {
	var total float64
	orm.NewOrm().Raw("SELECT COUNT(*) from storage_group_location where warehouse_location_id = ? and storage_group_id != ?", wl.ID, exID).QueryRow(&total)

	return total == 0
}
