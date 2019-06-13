// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"errors"
	"fmt"
	"strconv"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/delivery/model"
	"git.qasico.com/gudang/api/src/inventory"
	model3 "git.qasico.com/gudang/api/src/inventory/model"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	model4 "git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/warehouse"
)

var (
	errRequiredPartner             = "Tujuan warehouse harus diisi"
	errRequiredDocumentCode        = "Kode dokumen harus diisi"
	errRequiredPreparationPlan     = "Dokumen harus dipilih"
	errRequiredPreparationLocation = "Lokasi preparation harus diisi"
	errRequiredItems               = "Item harus diisi"
	errRequiredStockUnit           = "Unit harus dipilih"
	errRequiredQuantity            = "Quantity harus diisi"
	errInvalidCheckedBy            = "User tidak valid"
	errInvalidPreparation          = "Preparation dokumen tidak valid"
	errInvalidPartner              = "Asal warehouse tidak valid"
	errInvalidItemCode             = "Kode Item tidak valid"
	errInvalidBatchCode            = "Kode Batch tidak valid"
	errInvalidUnitCode             = "Kode Unit tidak valid"
	errInvalidPreparationDocument  = "Data tidak valid"
	errInvalidPreparationPlan      = "Dokumen tidak valid"
	errInvalidPreparationActual    = "Item tidak valid"
	errInvalidPreparationProgress  = "Masih ada unit yang dalam proses"
	errInvalidDate                 = "Tanggal tidak valid"
	errInvalidPreparationLocation  = "Lokasi preparation tidak valid"
	errInvalidStockUnit            = "Unit tidak valid"
	errInvalidStockUnitMovement    = "unit sedang ada proses movement, silahkan refresh browser"
	errInvalidStockUnitOpname      = "pada lokasi unit sedang ada proses stockopname"
	errInvalidLocation             = "pada lokasi sedang ada proses movement"
	errInvalidStockUnitDuplicate   = "unit ini sudah ada pada preparation unit"
	errInvalidQuantity             = "Quantity tidak mencukupi"
	errInvalidQuantityOver         = "Quantity melebihi dari yang dibutuhkan"
	errUniqueCode                  = "Kode dokumen telah digunakan"
)

func validPreparation(id int64, status string) (r *model.Preparation, e error) {
	e = orm.NewOrm().Raw("SELECT * FROM preparation WHERE id = ? and status = ?", id, status).QueryRow(&r)

	return
}

func validDocumentCode(code string, id int64) bool {
	var total float64
	orm.NewOrm().Raw("SELECT count(*) from preparation WHERE document_code = ? and id != ?", code, id).QueryRow(&total)

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

func validPreparationDocument(ide string) (rp *model.PreparationDocument, e error) {
	rp = new(model.PreparationDocument)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}

func validPreparationActual(ide string) (rp *model.PreparationActual, e error) {
	rp = new(model.PreparationActual)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}

func validSyncID(id int64) (r *model.Preparation, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT * FROM preparation WHERE id = ? and status = ?", id, "active").QueryRow(&r)

	var total float64
	o.Raw("SELECT count(*) FROM preparation_document WHERE preparation_id = ?", id).QueryRow(&total)

	if total != 0 {
		e = errors.New("already has item plan")
	}

	return
}

func validPreparationPlan(ide string) (rp *model.PreparationPlan, e error) {
	rp = new(model.PreparationPlan)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		rp.Status = "pending"
		e = rp.Read("id", "status")
	}

	return
}

func validFinishPreparation(r *model.Preparation) bool {
	o := orm.NewOrm()

	// harus ada plan
	var totalDocument float64
	o.Raw("SELECT count(*) FROM preparation_document where preparation_id = ?", r.ID).QueryRow(&totalDocument)
	if totalDocument == 0 {
		return false
	}

	// harus ada units
	// dan unit harus is_active semua
	var totalUnits float64
	o.Raw("SELECT count(*) FROM preparation_unit where preparation_id = ?", r.ID).QueryRow(&totalUnits)
	if totalUnits == 0 {
		return false
	}

	var totalPendingUnits float64
	o.Raw("SELECT count(*) FROM preparation_unit where preparation_id = ? and is_active = ?", r.ID, 1).QueryRow(&totalPendingUnits)

	return totalPendingUnits == totalUnits
}

func validPreparationLocation(ide string) (rp *warehouse.Location, e error) {
	rp = new(warehouse.Location)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		e = rp.Read()
	}

	return
}

func validStockUnit(ide string) (rp *model4.StockUnit, e error) {
	rp = new(model4.StockUnit)
	if rp.ID, e = common.Decrypt(ide); e == nil {
		rp.Status = "stored"
		if e = rp.Read("id", "status"); e == nil {
			e = rp.Storage.Read("ID")
		}
	}

	return
}

func validLocation(locID int64) (status bool) {
	var totalMove float64
	orm.NewOrm().Raw("SELECT count(*) FROM stock_movement where origin_id = ? AND status IN ('new','start');", locID).QueryRow(&totalMove)

	if totalMove > 0 {
		status = true
	}

	return
}

func validUnitStockopname(LocID int64) (status bool) {
	var OpnameID int64
	orm.NewOrm().Raw("SELECT so.id FROM stock_opname so "+
		"inner join warehouse_location wl on wl.id = so.location_id "+
		"where wl.id = ? and so.status = ?", LocID, "active").QueryRow(&OpnameID)
	if OpnameID > int64(0) {
		status = true
	}
	return
}

func validDuplicateUnitInPreparation(preparationID int64, unit int64) (status bool) {
	var totalPreUnit float64
	orm.NewOrm().Raw("SELECT count(*) FROM preparation_unit pu "+
		"inner join preparation p on p.id = pu.preparation_id "+
		"where p.id = ? and pu.unit_id = ? and p.status = ?", preparationID, unit, "active").QueryRow(&totalPreUnit)
	if totalPreUnit > 0 {
		status = true
	}
	return
}
