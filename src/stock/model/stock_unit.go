// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/custom/barcode"
	"git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/stock/services/unit"
	"git.qasico.com/gudang/api/src/user"
	"time"
)

func init() {
	orm.RegisterModel(new(StockUnit))
}

// StockUnit model for stock_unit table.
type StockUnit struct {
	ID           int64            `orm:"column(id);auto" json:"-"`
	Item         *model.Item      `orm:"column(item_id);null;rel(fk)" json:"item,omitempty"`
	Batch        *model.ItemBatch `orm:"column(batch_id);null;rel(fk)" json:"batch,omitempty"`
	Storage      *StockStorage    `orm:"column(storage_id);null;rel(fk)" json:"storage,omitempty"`
	RefID        uint64           `orm:"column(ref_id);null" json:"ref_id"`
	Code         string           `orm:"column(code);size(45)" json:"code"`
	Stock        float64          `orm:"column(stock);digits(20);decimals(2)" json:"stock"`
	IsDefect     int8             `orm:"column(is_defect)" json:"is_defect"`
	Status       string           `orm:"column(status);options(draft,stored,moving,prepared,void,out)" json:"status"`
	BarcodeImage string           `orm:"column(barcode_image);null" json:"barcode_image"`
	CreatedBy    *user.User       `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	ReceivedAt   time.Time        `orm:"column(received_at);type(timestamp)" json:"received_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *StockUnit) MarshalJSON() ([]byte, error) {
	type Alias StockUnit

	alias := &struct {
		ID          string `json:"id"`
		ItemID      string `json:"item_id"`
		BatchID     string `json:"batch_id"`
		StorageID   string `json:"storage_id"`
		CreatedByID string `json:"created_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	if m.Item != nil && m.Item.ID != int64(0) {
		alias.ItemID = common.Encrypt(m.Item.ID)
	} else {
		alias.Item = nil
	}

	if m.Batch != nil && m.Batch.ID != int64(0) {
		alias.BatchID = common.Encrypt(m.Batch.ID)
	} else {
		alias.Batch = nil
	}

	if m.Storage != nil && m.Storage.ID != int64(0) {
		alias.StorageID = common.Encrypt(m.Storage.ID)
	} else {
		alias.Storage = nil
	}

	if m.CreatedBy != nil && m.CreatedBy.ID != int64(0) {
		alias.CreatedByID = common.Encrypt(m.CreatedBy.ID)
	} else {
		alias.CreatedBy = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating StockUnit struct into stock_unit table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to stock_unit.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *StockUnit) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting stock_unit data
// this also will truncated all data from all table
// that have relation with this stock_unit.
func (m *StockUnit) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}

// Read execute select based on data struct that already assigned.
func (m *StockUnit) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}

// GenerateCode function to make new code and also barcode image
// it is apply for new stock unit only
func (m *StockUnit) GenerateCode(parent string) {
	if m.Code == "" {
		if parent == "" {
			m.Code = unit.NewCode()
		} else {
			m.Code = unit.NewChildCode(parent)
		}
	}

	m.BarcodeImage, _ = barcode.MakeBarcode(m.Code)
}

func (m *StockUnit) SetStored(storage *StockStorage) error {
	m.Status = "stored"

	if storage == nil {
		return m.Save("status")
	}

	m.Storage = storage

	return m.Save("status", "storage")
}

func (m *StockUnit) SetMoving() error {
	m.Status = "moving"

	return m.Save("status")
}

func (m *StockUnit) SetPrepared() error {
	m.Status = "prepared"

	return m.Save("status")
}

func (m *StockUnit) SetVoid() error {
	m.Status = "void"
	m.Storage = nil

	return m.Save("status", "storage")
}

func (m *StockUnit) SetOut() error {
	m.Status = "out"

	return m.Save("status")
}
