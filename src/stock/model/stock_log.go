// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
)

func init() {
	orm.RegisterModel(new(StockLog))
}

// StockLog model for stock_log table.
type StockLog struct {
	ID         int64            `orm:"column(id);auto" json:"-"`
	StockUnit  *StockUnit       `orm:"column(stock_unit_id);null;rel(fk)" json:"stock_unit,omitempty"`
	Item       *model.Item      `orm:"column(item_id);null;rel(fk)" json:"item,omitempty"`
	Batch      *model.ItemBatch `orm:"column(batch_id);null;rel(fk)" json:"batch,omitempty"`
	RefType    string           `orm:"column(ref_type);size(45);null" json:"ref_type"`
	RefID      uint64           `orm:"column(ref_id);null" json:"-"`
	RefCode    string           `orm:"column(ref_code);null" json:"ref_code"`
	Quantity   float64          `orm:"column(quantity);null;digits(20);decimals(2)" json:"quantity"`
	RecordedAt time.Time        `orm:"column(recorded_at);type(timestamp)" json:"recorded_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *StockLog) MarshalJSON() ([]byte, error) {
	type Alias StockLog

	alias := &struct {
		ID          string `json:"id"`
		StockUnitID string `json:"stock_unit_id"`
		ItemID      string `json:"item_id"`
		BatchID     string `json:"batch_id"`
		RefID       string `json:"ref_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.StockUnitID when m.StockUnit not nill
	// and the ID is setted
	if m.StockUnit != nil && m.StockUnit.ID != int64(0) {
		alias.StockUnitID = common.Encrypt(m.StockUnit.ID)
	} else {
		alias.StockUnit = nil
	}

	// Encrypt alias.ItemID when m.Item not nill
	// and the ID is setted
	if m.Item != nil && m.Item.ID != int64(0) {
		alias.ItemID = common.Encrypt(m.Item.ID)
	} else {
		alias.Item = nil
	}

	if m.RefID > 0 {
		alias.RefID = common.Encrypt(m.RefID)
	}

	// Encrypt alias.BatchID when m.Batch not nill
	// and the ID is setted
	if m.Batch != nil && m.Batch.ID != int64(0) {
		alias.BatchID = common.Encrypt(m.Batch.ID)
	} else {
		alias.Batch = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating StockLog struct into stock_log table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to stock_log.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *StockLog) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting stock_log data
// this also will truncated all data from all table
// that have relation with this stock_log.
func (m *StockLog) Delete() (err error) {
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

// Read execute select based on data struct that already
// assigned.
func (m *StockLog) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
