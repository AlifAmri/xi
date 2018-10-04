// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/receiving/model"
)

var (
	errRequiredDocumentCode = "Kode dokumen harus diisi"
	errRequiredReceivingAt  = "Tanggal harus diisi"
	errRequiredItems        = "Item harus diisi"
	errInvalidReceiptPlan   = "Dokumen tidak valid"
	errInvalidDate          = "Tanggal tidak valid"
	errInvalidPartner       = "Asal warehouse tidak valid"
	errInvalidReceivingPlan = "Data tidak valid"
	errUniqueCode           = "Kode dokumen telah digunakan"
)

func validReceiptPlan(id int64, status string) (r *model.ReceiptPlan, e error) {
	e = orm.NewOrm().Raw("SELECT * FROM receipt_plan WHERE id = ? and status = ?", id, status).QueryRow(&r)

	return
}

func validDocumentCode(code string, id int64) bool {
	var total float64
	orm.NewOrm().Raw("SELECT count(*) from receipt_plan WHERE document_code = ? and id != ?", code, id).QueryRow(&total)

	return total == 0
}

func validPartner(ide string) (rp *model2.Partnership, e error) {
	rp = new(model2.Partnership)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}

func validReceiptPlanItem(ide string) (rp *model.ReceiptPlanItem, e error) {
	rp = new(model.ReceiptPlanItem)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}
