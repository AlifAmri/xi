// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"fmt"

	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Item))
}

// Item model for item table.
type Item struct {
	ID            int64           `orm:"column(id);auto" json:"-"`
	Group         *ItemGroup      `orm:"column(group_id);null;rel(fk)" json:"group,omitempty"`
	Type          *ItemType       `orm:"column(type_id);null;rel(fk)" json:"type,omitempty"`
	Category      *ItemCategory   `orm:"column(category_id);null;rel(fk)" json:"category,omitempty"`
	PreferredArea *warehouse.Area `orm:"column(preferred_area);null;rel(fk)" json:"preferred_area"`
	Code          string          `orm:"column(code);size(45)" json:"code"`
	IsActive      int8            `orm:"column(is_active)" json:"is_active"`
	Name          string          `orm:"column(name);size(145)" json:"name"`
	Image         string          `orm:"column(image);null" json:"image"`
	Stock         float64         `orm:"column(stock)" json:"stock"`
	BarcodeNumber string          `orm:"column(barcode_number);size(50);null" json:"barcode_number"`
	BarcodeImage  string          `orm:"column(barcode_image);null" json:"barcode_image"`
	Note          string          `orm:"column(note);null" json:"note"`

	Attributes map[string]string `orm:"-" json:"attributes"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Item) MarshalJSON() ([]byte, error) {
	type Alias Item

	alias := &struct {
		ID              string `json:"id"`
		TypeID          string `json:"type_id"`
		CategoryID      string `json:"category_id"`
		PreferredAreaID string `json:"preferred_area_id"`
		GroupID         string `json:"group_id"`
		AliasName       string `json:"alias"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.GroupID when m.Group not nill
	// and the ID is setted
	if m.Group != nil && m.Group.ID != int64(0) {
		alias.GroupID = common.Encrypt(m.Group.ID)
	} else {
		alias.Group = nil
	}

	// Encrypt alias.TypeID when m.Type not nill
	// and the ID is setted
	if m.Type != nil && m.Type.ID != int64(0) {
		alias.TypeID = common.Encrypt(m.Type.ID)
	} else {
		alias.Type = nil
	}

	// Encrypt alias.CategoryID when m.Category not nill
	// and the ID is setted
	if m.Category != nil && m.Category.ID != int64(0) {
		alias.CategoryID = common.Encrypt(m.Category.ID)
	} else {
		alias.Category = nil
	}

	// Encrypt alias.PreferredAreaID when m.PreferredArea not nill
	// and the ID is setted
	if m.PreferredArea != nil && m.PreferredArea.ID != int64(0) {
		alias.PreferredAreaID = common.Encrypt(m.PreferredArea.ID)
	} else {
		alias.PreferredArea = nil
	}

	alias.AliasName = fmt.Sprintf("%s - %s", m.Code, m.Name)

	return json.Marshal(alias)
}

// Save inserting or updating Item struct into item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Item) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting item data
// this also will truncated all data from all table
// that have relation with this item.
func (m *Item) Delete() (err error) {
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
func (m *Item) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
