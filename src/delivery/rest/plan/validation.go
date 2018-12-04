// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"errors"
	"fmt"
	"strconv"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
)

var (
	errRequiredDocumentCode = "Kode dokumen harus diisi"
	errRequiredReceivingAt  = "Tanggal harus diisi"
	errRequiredItems        = "Item harus diisi"
	errInvalidReceiptPlan   = "Dokumen tidak valid"
	errInvalidDate          = "Tanggal tidak valid"
	errInvalidPartner       = "Tujuan pengiriman tidak valid"
	errInvalidReceivingPlan = "Data tidak valid"
	errInvalidBatchCode     = "Kode batch tidak valid"
	errUniqueCode           = "Kode dokumen telah digunakan"
)

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

func validPreparationPlan(id int64, status string) (r *model.PreparationPlan, e error) {
	e = orm.NewOrm().Raw("SELECT * FROM preparation_plan WHERE id = ? and status = ?", id, status).QueryRow(&r)

	return
}

func validDocumentCode(code string, id int64) bool {
	var total float64
	orm.NewOrm().Raw("SELECT count(*) from preparation_plan WHERE document_code = ? and id != ?", code, id).QueryRow(&total)

	return total == 0
}

func validPartner(ide string) (rp *model2.Partnership, e error) {
	rp = new(model2.Partnership)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}

func validPreparationPlanItem(ide string) (rp *model.PreparationPlanItem, e error) {
	rp = new(model.PreparationPlanItem)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}
