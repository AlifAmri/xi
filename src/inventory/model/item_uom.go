// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(ItemUom))
}

// ItemUom model for item_uom table.
type ItemUom struct {
	ID   int64     `orm:"column(id);auto" json:"-"`
	Type *ItemType `orm:"column(type_id);null;rel(fk)" json:"type,omitempty"`
	Name string    `orm:"column(name);size(45)" json:"name"`
	Note string    `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ItemUom) MarshalJSON() ([]byte, error) {
	type Alias ItemUom

	alias := &struct {
		ID     string `json:"id"`
		TypeID string `json:"type_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.TypeID when m.Type not nill
	// and the ID is setted
	if m.Type != nil && m.Type.ID != int64(0) {
		alias.TypeID = common.Encrypt(m.Type.ID)
	} else {
		alias.Type = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating ItemUom struct into item_uom table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to item_uom.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ItemUom) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting item_uom data
// this also will truncated all data from all table
// that have relation with this item_uom.
func (m *ItemUom) Delete() (err error) {
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
func (m *ItemUom) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
