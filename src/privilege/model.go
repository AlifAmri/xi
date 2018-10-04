// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package privilege

import (
	"encoding/json"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Privilege))
}

// Privilege model for privilege table.
type Privilege struct {
	ID       int64  `orm:"column(id);auto" json:"-"`
	Name     string `orm:"column(name);size(45)" json:"name"`
	Action   string `orm:"column(action);size(45)" json:"action"`
	IsActive uint8  `orm:"column(is_active);null" json:"is_active"`
	Note     string `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Privilege) MarshalJSON() ([]byte, error) {
	type Alias Privilege

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save inserting or updating Privilege struct into privilege table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to privilege.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Privilege) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting privilege data
// this also will truncated all data from all table
// that have relation with this privilege.
func (m *Privilege) Delete() (err error) {
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
func (m *Privilege) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}

// Activate mengupdate is active menjadi 1
func (m *Privilege) Activate() (e error) {
	if e = m.Read("id"); e == nil {
		m.IsActive = 1
		e = m.Save("is_active")
	}

	return
}

// Deactivate mengupdate is active menjadi 0
func (m *Privilege) Deactivate() (e error) {
	if e = m.Read("id"); e == nil {
		m.IsActive = 0
		e = m.Save("is_active")
	}

	return
}
