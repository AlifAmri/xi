// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"errors"
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
)

var (
	errRequiredItems            = "Item harus diisi"
	errInvalidDeliveryOrder     = "Surat Jalan tidak valid"
	errInvalidDeliveryOrderItem = "Dokumen tidak valid"
	errInvalidPartner           = "Partner tidak valid"
	errInvalidPreparation       = "Data tidak valid"
)

func validDeliveryOrder(id int64, status string) (r *model.DeliveryOrder, e error) {
	e = orm.NewOrm().Raw("SELECT * FROM delivery_order WHERE id = ? and status = ?", id, status).QueryRow(&r)

	return
}

func validPartner(ide string) (rp *model2.Partnership, e error) {
	rp = new(model2.Partnership)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}

func validPreparation(ide string, do *model.DeliveryOrder) (rp *model.Preparation, e error) {
	rp = new(model.Preparation)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		if e = rp.Read(); e == nil {
			if rp.DeliveryOrder != nil && rp.DeliveryOrder.ID != do.ID {
				e = errors.New("preparation sudah ada surat jalan")
			}
		}
	}

	return
}
