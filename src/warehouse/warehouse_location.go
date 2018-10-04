// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package warehouse

import (
	"encoding/json"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Location))
}

// Location model for warehouse_location table.
type Location struct {
	ID              int64  `orm:"column(id);auto" json:"-"`
	Area            *Area  `orm:"column(warehouse_area_id);rel(fk)" json:"warehouse_area,omitempty"`
	Code            string `orm:"column(code);size(45)" json:"code"`
	Name            string `orm:"column(name);size(45)" json:"name"`
	IsActive        int8   `orm:"column(is_active);null" json:"is_active"`
	CoordinateX     int    `orm:"column(coordinate_x)" json:"coordinate_x"`
	CoordinateY     int    `orm:"column(coordinate_y)" json:"coordinate_y"`
	CoordinateW     int    `orm:"column(coordinate_w)" json:"coordinate_w"`
	CoordinateH     int    `orm:"column(coordinate_h)" json:"coordinate_h"`
	StorageCapacity int    `orm:"column(storage_capacity)" json:"storage_capacity"`
	StorageUsed     int    `orm:"column(storage_used)" json:"storage_used"`
	Note            string `orm:"column(note);null" json:"note"`
}

// TableName custom table name
func (m *Location) TableName() string {
	return "warehouse_location"
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Location) MarshalJSON() ([]byte, error) {
	type Alias Location

	alias := &struct {
		ID     string `json:"id"`
		AreaID string `json:"warehouse_area_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.AreaID when m.Area not nill
	// and the ID is setted
	if m.Area != nil && m.Area.ID != int64(0) {
		alias.AreaID = common.Encrypt(m.Area.ID)
	} else {
		alias.Area = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating Location struct into warehouse_location table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to warehouse_location.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Location) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting warehouse_location data
// this also will truncated all data from all table
// that have relation with this warehouse_location.
func (m *Location) Delete() (err error) {
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
func (m *Location) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
