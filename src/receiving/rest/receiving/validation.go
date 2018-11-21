// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"errors"
	"fmt"
	"strconv"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory"
	model3 "git.qasico.com/gudang/api/src/inventory/model"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/receiving/model"
	model4 "git.qasico.com/gudang/api/src/stock/model"
)

var (
	errRequiredPartner          = "Asal warehouse harus diisi"
	errRequiredDocumentCode     = "Kode dokumen harus diisi"
	errRequiredReceiptPlan      = "Dokumen harus dipilih"
	errRequiredItems            = "Item harus diisi"
	errInvalidCheckedBy         = "User tidak valid"
	errInvalidReceiving         = "Receiving dokumen tidak valid"
	errInvalidPartner           = "Asal warehouse tidak valid"
	errInvalidItemCode          = "Kode Item tidak valid"
	errInvalidBatchCode         = "Kode Batch tidak valid"
	errInvalidUnitCode          = "Kode Unit tidak valid"
	errInvalidReceivingPlan     = "Data tidak valid"
	errInvalidReceiptPlan       = "Dokumen tidak valid"
	errInvalidReceivingActual   = "Item tidak valid"
	errInvalidReceivingProgress = "Masih ada unit yang dalam proses atau belum ada surat jalan"
	errUniqueCode               = "Kode dokumen telah digunakan"
)

func validReceiving(id int64, status string) (r *model.Receiving, e error) {
	e = orm.NewOrm().Raw("SELECT * FROM receiving WHERE id = ? and status = ?", id, status).QueryRow(&r)

	return
}

func validDocumentCode(code string, id int64) bool {
	var total float64
	orm.NewOrm().Raw("SELECT count(*) from receiving WHERE document_code = ? and id != ?", code, id).QueryRow(&total)

	return total == 0
}

func validPartner(ide string) (rp *model2.Partnership, e error) {
	rp = new(model2.Partnership)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}

func validItemCode(code string) (i *model3.Item, new int8, e error) {
	i = &model3.Item{Code: code, Type: &model3.ItemType{ID: 1}}

	if e = i.Read("type_id", "code"); e != nil {
		e = i.Save()
		new = 1
	}

	return
}

func validBatchCode(code string, i *model3.Item) (b *model3.ItemBatch, e error) {
	b = inventory.GetBatch(i.ID, code)

	return
}

func validBatchCodeString(code string) (c string, e error) {
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

func validStockUnit(code string, i *model3.Item, b *model3.ItemBatch) (s *model4.StockUnit, e error) {
	orm.NewOrm().Raw("SELECT * FROM stock_unit where code = ?", code).QueryRow(&s)

	if s == nil {
		s = &model4.StockUnit{Item: i, Batch: b, Code: code, Status: "draft"}
		s.Save()
	}

	return
}

func validReceivingPlan(ide string) (rp *model.ReceivingDocument, e error) {
	rp = new(model.ReceivingDocument)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}

func validReceivingActual(ide string) (rp *model.ReceivingActual, e error) {
	rp = new(model.ReceivingActual)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}

func validSyncID(id int64) (r *model.Receiving, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT * FROM receiving WHERE id = ? and status = ?", id, "active").QueryRow(&r)

	var total float64
	o.Raw("SELECT count(*) FROM receiving_document WHERE receiving_id = ?", id).QueryRow(&total)

	if total != 0 {
		e = errors.New("already has item plan")
	}

	return
}

func validReceiptPlan(ide string) (rp *model.ReceiptPlan, e error) {
	rp = new(model.ReceiptPlan)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		rp.Status = "pending"
		e = rp.Read("id", "status")
	}

	return
}

func validFinishReceiving(r *model.Receiving) bool {
	o := orm.NewOrm()

	// harus ada plan
	var totalDocument float64
	o.Raw("SELECT count(*) FROM receiving_document where receiving_id = ?", r.ID).QueryRow(&totalDocument)
	if totalDocument == 0 {
		return false
	}

	// harus ada units
	// dan unit harus is_active semua
	var totalUnits float64
	o.Raw("SELECT count(*) FROM receiving_unit where receiving_id = ?", r.ID).QueryRow(&totalUnits)
	if totalUnits == 0 {
		return false
	}

	var totalPendingUnits float64
	o.Raw("SELECT count(*) FROM receiving_unit where receiving_id = ? and is_active = ?", r.ID, 1).QueryRow(&totalPendingUnits)

	return totalPendingUnits == totalUnits
}
