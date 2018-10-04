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
	orm.RegisterModel(new(StorageGroup))
}

// StorageArea model for warehouse_area table.
type StorageArea struct {
	Identity   int64  `json:"-"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Note       string `json:"note"`
	IsSelected int8   `json:"is_selected"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *StorageArea) MarshalJSON() ([]byte, error) {
	type Alias StorageArea

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.Identity),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// StorageGroup model for storage_group table.
type StorageGroup struct {
	ID        int64  `orm:"column(id);pk" json:"-"`
	Name      string `orm:"column(name);size(45)" json:"name"`
	Type      string `orm:"column(type);options(default,item_group,item_category,item_batch,item_code,ncp)" json:"type"`
	TypeValue string `orm:"column(type_value);size(45);null" json:"type_value"`
	IsActive  int8   `orm:"column(is_active)" json:"is_active"`
	IsPrimary int8   `orm:"column(is_primary)" json:"is_primary"`
	Note      string `orm:"column(note);null" json:"note"`

	Areas        []*StorageArea     `orm:"-" json:"areas"`
	ItemGroup    model.ItemGroup    `orm:"-" json:"item_group,omitempty"`
	ItemCategory model.ItemCategory `orm:"-" json:"item_category,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *StorageGroup) MarshalJSON() ([]byte, error) {
	type Alias StorageGroup

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save inserting or updating StorageGroup struct into storage_group table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to storage_group.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *StorageGroup) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting storage_group data
// this also will truncated all data from all table
// that have relation with this storage_group.
func (m *StorageGroup) Delete() (err error) {
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
func (m *StorageGroup) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
