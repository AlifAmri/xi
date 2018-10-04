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
	orm.RegisterModel(new(Warehouse))
}

// Warehouse model for warehouse table.
type Warehouse struct {
	ID   int64  `orm:"column(id);auto" json:"-"`
	Code string `orm:"column(code);size(45)" json:"code"`
	Name string `orm:"column(name);size(45)" json:"name"`
	Note string `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Warehouse) MarshalJSON() ([]byte, error) {
	type Alias Warehouse

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save inserting or updating Warehouse struct into warehouse table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to warehouse.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Warehouse) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting warehouse data
// this also will truncated all data from all table
// that have relation with this warehouse.
func (m *Warehouse) Delete() (err error) {
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
func (m *Warehouse) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
