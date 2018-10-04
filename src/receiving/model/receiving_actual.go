// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.qasico.com/gudang/api/src/inventory/model"
	stock "git.qasico.com/gudang/api/src/stock/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(ReceivingActual))
}

// ReceivingActual model for receiving_actual table.
type ReceivingActual struct {
	ID               int64            `orm:"column(id);auto" json:"-"`
	Receiving        *Receiving       `orm:"column(receiving_id);rel(fk)" json:"receiving,omitempty"`
	Item             *model.Item      `orm:"column(item_id);rel(fk)" json:"item,omitempty"`
	Batch            *model.ItemBatch `orm:"column(batch_id);null;rel(fk)" json:"batch,omitempty"`
	Unit             *stock.StockUnit `orm:"column(unit_id);null;rel(fk)" json:"unit,omitempty"`
	QuantityPlanned  float64          `orm:"column(quantity_planned);digits(12);decimals(2)" json:"quantity_planned"`
	QuantityReceived float64          `orm:"column(quantity_received);digits(12);decimals(2)" json:"quantity_received"`
	QuantityDefect   float64          `orm:"column(quantity_defect);digits(12);decimals(2)" json:"quantity_defect"`
	Note             string           `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ReceivingActual) MarshalJSON() ([]byte, error) {
	type Alias ReceivingActual

	alias := &struct {
		ID          string `json:"id"`
		ReceivingID string `json:"receiving_id"`
		ItemID      string `json:"item_id"`
		BatchID     string `json:"batch_id"`
		UnitID      string `json:"unit_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.BatchID when m.Batch not nill
	// and the ID is setted
	if m.Batch != nil && m.Batch.ID != int64(0) {
		alias.BatchID = common.Encrypt(m.Batch.ID)
	} else {
		alias.Batch = nil
	}

	// Encrypt alias.UnitID when m.Unit not nill
	// and the ID is setted
	if m.Unit != nil && m.Unit.ID != int64(0) {
		alias.UnitID = common.Encrypt(m.Unit.ID)
	} else {
		alias.Unit = nil
	}

	// Encrypt alias.ReceivingID when m.Receiving not nill
	// and the ID is setted
	if m.Receiving != nil && m.Receiving.ID != int64(0) {
		alias.ReceivingID = common.Encrypt(m.Receiving.ID)
	} else {
		alias.Receiving = nil
	}

	// Encrypt alias.ItemID when m.Item not nill
	// and the ID is setted
	if m.Item != nil && m.Item.ID != int64(0) {
		alias.ItemID = common.Encrypt(m.Item.ID)
	} else {
		alias.Item = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating ReceivingActual struct into receiving_actual table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to receiving_actual.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ReceivingActual) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting receiving_actual data
// this also will truncated all data from all table
// that have relation with this receiving_actual.
func (m *ReceivingActual) Delete() (err error) {
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
func (m *ReceivingActual) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
