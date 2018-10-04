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
	orm.RegisterModel(new(Partnership))
}

// Partnership model for partnership table.
type Partnership struct {
	ID             int64            `orm:"column(id);auto" json:"-"`
	Type           *PartnershipType `orm:"column(type_id);rel(fk)" json:"type,omitempty"`
	CompanyName    string           `orm:"column(company_name);size(100);null" json:"company_name"`
	CompanyAddress string           `orm:"column(company_address);null" json:"company_address"`
	CompanyPhone   string           `orm:"column(company_phone);size(45);null" json:"company_phone"`
	CompanyEmail   string           `orm:"column(company_email);size(100);null" json:"company_email"`
	ContactPerson  string           `orm:"column(contact_person);size(100);null" json:"contact_person"`
	IsActive       int8             `orm:"column(is_active)" json:"is_active"`
	Note           string           `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Partnership) MarshalJSON() ([]byte, error) {
	type Alias Partnership

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

// Save inserting or updating Partnership struct into partnership table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to partnership.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Partnership) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting partnership data
// this also will truncated all data from all table
// that have relation with this partnership.
func (m *Partnership) Delete() (err error) {
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
func (m *Partnership) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
