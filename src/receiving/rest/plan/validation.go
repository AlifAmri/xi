// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"errors"
	"fmt"
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/receiving/model"
	"strconv"
)

var (
	errRequiredDocumentCode = "Kode dokumen harus diisi"
	errRequiredReceivingAt  = "Tanggal harus diisi"
	errRequiredItems        = "Item harus diisi"
	errInvalidReceiptPlan   = "Dokumen tidak valid"
	errInvalidDate          = "Tanggal tidak valid"
	errInvalidPartner       = "Asal warehouse tidak valid"
	errInvalidReceivingPlan = "Data tidak valid"
	errInvalidUnit          = "Kode stock sudah dipakai"
	errUniqueUnit           = "Kode stock harus unik"
	errInvalidBatchCode     = "Kode batch tidak valid"
	errUniqueCode           = "Kode dokumen telah digunakan"
)

func validUnitCode(code string, exID string) bool {
	var total int64
	o := orm.NewOrm()

	o.Raw("SELECT count(*) FROM stock_unit where code = ?", code).QueryRow(&total)
	if total == 0 {
		id, _ := common.Decrypt(exID)
		o.Raw("SELECT count(*) from receipt_plan_item where unit_code = ? and id != ?", code, id).QueryRow(&total)
	}

	return total == 0
}

func validBatchCode(code string) (c string, e error) {
	c = code

	if len(c) == 3 {
		c = fmt.Sprintf("%s%s", "0", code)
	}

	if len(c) != 4 {
		e = errors.New("wrong format")
	} else {
		cx := c[0:2]
		if !validWeek(cx) {
			return "", errors.New("not valid")
		}
	}

	return
}

func validWeek(s string) bool {
	i, e := strconv.Atoi(s)
	if e == nil && i > 0 && i < 54 {
		return true
	}

	return false
}

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
