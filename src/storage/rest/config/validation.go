// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/storage/model"
)

var (
	errRequiredName        = "nama harus diisi"
	errRequiredType        = "tipe harus dipilih"
	errRequiredValue       = "value harus diisi"
	errUniqueName          = "nama tersebut telah digunakan"
	errInvalidStorageGroup = "grup storage tidak dapat ditemukan"
	errAlreadyActived      = "grup storage sudah aktif"
	errAlreadyDeactived    = "grup storage sudah tidak aktif"
	errRequiredArea        = "Area harus dipilih"
)

type area struct {
	AreaID          string `json:"id"`
	FromLocationID  string `json:"form_location_id"`
	UntilLocationID string `json:"until_location_id"`
}

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
