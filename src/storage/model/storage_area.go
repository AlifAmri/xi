// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/warehouse"
)

func init() {
	orm.RegisterModel(new(StorageGroupArea))
}

// StorageGroupArea model for StorageGroupArea table.
type StorageGroupArea struct {
	ID            int64               `orm:"column(id);auto" json:"-"`
	StorageGroup  *StorageGroup       `orm:"column(storage_group_id);rel(fk)" json:"storage_group,omitempty"`
	WarehouseArea *warehouse.Area     `orm:"column(warehouse_area_id);rel(fk)" json:"warehouse_area,omitempty"`
	LocationFrom  *warehouse.Location `orm:"column(location_from_id);rel(fk)" json:"location_from,omitempty"`
	LocationEnd   *warehouse.Location `orm:"column(location_end_id);rel(fk)" json:"location_end,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *StorageGroupArea) MarshalJSON() ([]byte, error) {
	type Alias StorageGroupArea

	alias := &struct {
		ID              string `json:"id"`
		WarehouseAreaID string `json:"warehouse_area_id"`
		LocationFromID  string `json:"location_from_id"`
		LocationEndID   string `json:"location_end_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	if m.WarehouseArea != nil && m.WarehouseArea.ID != int64(0) {
		alias.WarehouseAreaID = common.Encrypt(m.WarehouseArea.ID)
	} else {
		alias.WarehouseArea = nil
	}

	if m.LocationFrom != nil && m.LocationFrom.ID != int64(0) {
		alias.LocationFromID = common.Encrypt(m.LocationFrom.ID)
	} else {
		alias.LocationFrom = nil
	}

	if m.LocationEnd != nil && m.LocationEnd.ID != int64(0) {
		alias.LocationEndID = common.Encrypt(m.LocationEnd.ID)
	} else {
		alias.LocationEnd = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating StorageGroupArea struct into StorageGroupArea table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to StorageGroupArea.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *StorageGroupArea) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting StorageGroupArea data
// this also will truncated all data from all table
// that have relation with this StorageGroupArea.
func (m *StorageGroupArea) Delete() (err error) {
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
func (m *StorageGroupArea) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
