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
	orm.RegisterModel(new(Area))
}

// Area model for warehouse_area table.
type Area struct {
	ID        int64      `orm:"column(id);auto" json:"-"`
	Warehouse *Warehouse `orm:"column(warehouse_id);rel(fk)" json:"warehouse,omitempty"`
	Type      string     `orm:"column(type);options(storage,receiving,preparation,other)" json:"type"`
	Code      string     `orm:"column(code);size(45)" json:"code"`
	Name      string     `orm:"column(name);size(45)" json:"name"`
	IsActive  int8       `orm:"column(is_active);null" json:"is_active"`
	Note      string     `orm:"column(note);null" json:"note"`
}

// TableName custom table name
func (m *Area) TableName() string {
	return "warehouse_area"
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Area) MarshalJSON() ([]byte, error) {
	type Alias Area

	alias := &struct {
		ID          string `json:"id"`
		WarehouseID string `json:"warehouse_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.WarehouseID when m.Warehouse not nill
	// and the ID is setted
	if m.Warehouse != nil && m.Warehouse.ID != int64(0) {
		alias.WarehouseID = common.Encrypt(m.Warehouse.ID)
	} else {
		alias.Warehouse = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating Area struct into warehouse_area table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to warehouse_area.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Area) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting warehouse_area data
// this also will truncated all data from all table
// that have relation with this warehouse_area.
func (m *Area) Delete() (err error) {
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
func (m *Area) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
