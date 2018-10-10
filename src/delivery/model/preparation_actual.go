// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
)

func init() {
	orm.RegisterModel(new(PreparationActual))
}

// PreparationActual model for Preparation_actual table.
type PreparationActual struct {
	ID               int64            `orm:"column(id);auto" json:"-"`
	Preparation      *Preparation     `orm:"column(preparation_id);rel(fk)" json:"Preparation,omitempty"`
	Item             *model.Item      `orm:"column(item_id);rel(fk)" json:"item,omitempty"`
	Batch            *model.ItemBatch `orm:"column(batch_id);null;rel(fk)" json:"batch,omitempty"`
	QuantityPlanned  float64          `orm:"column(quantity_planned);digits(12);decimals(2)" json:"quantity_planned"`
	QuantityPrepared float64          `orm:"column(quantity_prepared);digits(12);decimals(2)" json:"quantity_prepared"`
	Note             string           `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PreparationActual) MarshalJSON() ([]byte, error) {
	type Alias PreparationActual

	alias := &struct {
		ID            string `json:"id"`
		PreparationID string `json:"preparation_id"`
		ItemID        string `json:"item_id"`
		BatchID       string `json:"batch_id"`
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

	// Encrypt alias.PreparationID when m.Preparation not nill
	// and the ID is setted
	if m.Preparation != nil && m.Preparation.ID != int64(0) {
		alias.PreparationID = common.Encrypt(m.Preparation.ID)
	} else {
		alias.Preparation = nil
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

// Save inserting or updating PreparationActual struct into Preparation_actual table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to Preparation_actual.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *PreparationActual) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting Preparation_actual data
// this also will truncated all data from all table
// that have relation with this Preparation_actual.
func (m *PreparationActual) Delete() (err error) {
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
func (m *PreparationActual) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
