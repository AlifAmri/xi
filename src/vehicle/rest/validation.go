// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	pmodel "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/vehicle/model"
)

var (
	errRequiredPurpose       = "tujuan kedatangan harus diisi"
	errRequiredVehicleType   = "tipe kendaraan harus dipilih"
	errRequiredVehicleNumber = "nomor kendaraan harus diisi"
	errRequiredDriver        = "nama driver harus diisi"
	errRequiredPicture       = "foto id driver harus diisi"
	errInvalidID             = "kedatangan kendaraan tidak ditemukan"
	errSubconInvalid         = "subcon tidak ditemukan"
)

func validID(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM incoming_vehicle where id = ? and status = ?", id, "in_progress").QueryRow(&total)

	return total == 1
}

func validVehicle(id int64) (i *model.IncomingVehicle, e error) {
	i = &model.IncomingVehicle{ID: id}
	e = i.Read()

	return
}

func validSubcon(ide string) (ug *pmodel.Partnership, e error) {
	ug = new(pmodel.Partnership)
	if ug.ID, e = common.Decrypt(ide); e == nil {
		e = ug.Read()
	}

	return
}
