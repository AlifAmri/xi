// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package partner

import (
	"git.qasico.com/gudang/api/src/partnership/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

var (
	errRequiredName     = "nama harus diisi"
	errRequiredType     = "type partnership harus dipilih"
	errUniqueName       = "nama tersebut telah digunakan"
	errInvalidID        = "partner tidak dapat ditemukan"
	errInvalidType      = "type tidak dapat ditemukan"
	errCascadeID        = "partner tersebut tidak dapat dihapus"
	errAlreadyDeactived = "status partner sudah tidak aktif"
	errAlreadyActived   = "status partner sudah aktif"
)

func validName(name string, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM partnership where company_name = ? and id != ?", name, exclude).QueryRow(&total)

	return total == 0
}

func validID(id int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM partnership where id = ?", id).QueryRow(&total)

	return total == 1
}

func validType(ide string) (it *model.PartnershipType, e error) {
	it = new(model.PartnershipType)
	if it.ID, e = common.Decrypt(ide); e == nil {
		e = it.Read()
	}

	return
}

func validPartner(id int64) (u *model.Partnership, e error) {
	u = &model.Partnership{ID: id}
	e = u.Read()

	return
}

func validDelete(id int64) bool {
	// @todo validasi kondisi partnerhip dapat dihapus
	return id == -1
}
