// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(ItemBatch))
}

// ItemBatch model for item_batch table.
type ItemBatch struct {
	ID        int64     `orm:"column(id);auto" json:"-"`
	Item      *Item     `orm:"column(item_id);rel(fk)" json:"item,omitempty"`
	Code      string    `orm:"column(code);size(45)" json:"code"`
	Name      string    `orm:"column(name);size(45)" json:"name"`
	Stock     float64   `orm:"column(stock)" json:"stock"`
	ProduceAt time.Time `orm:"column(produced_at);type(timestamp);null" json:"produced_at"`
	ExpiredAt time.Time `orm:"column(expired_at);type(timestamp);null" json:"expired_at"`
	EntryAt   time.Time `orm:"column(entry_at);type(timestamp);null" json:"entry_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ItemBatch) MarshalJSON() ([]byte, error) {
	type Alias ItemBatch

	alias := &struct {
		ID     string `json:"id"`
		ItemID string `json:"item_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
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

// Save inserting or updating ItemBatch struct into item_batch table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to item_batch.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ItemBatch) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting item_batch data
// this also will truncated all data from all table
// that have relation with this item_batch.
func (m *ItemBatch) Delete() (err error) {
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
func (m *ItemBatch) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
