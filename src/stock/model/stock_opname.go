// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/warehouse"
)

func init() {
	orm.RegisterModel(new(StockOpname))
}

// StockOpname model for stock_opname table.
type StockOpname struct {
	ID         int64               `orm:"column(id);auto" json:"-"`
	Code       string              `orm:"column(code);size(45)" json:"code"`
	Location   *warehouse.Location `orm:"column(location_id);null;rel(fk)" json:"location,omitempty"`
	Type       string              `orm:"column(type);options(opname,adjustment)" json:"type"`
	Status     string              `orm:"column(status);options(active,finish,cancelled)" json:"status"`
	Note       string              `orm:"column(note);null" json:"note"`
	CreatedBy  *user.User          `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	ApprovedBy *user.User          `orm:"column(approved_by);null;rel(fk)" json:"approved_by"`
	CreatedAt  time.Time           `orm:"column(created_at);type(timestamp)" json:"created_at"`
	ApprovedAt time.Time           `orm:"column(approved_at);type(timestamp);null" json:"approved_at"`

	Items []*StockOpnameItem `orm:"reverse(many)" json:"items,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *StockOpname) MarshalJSON() ([]byte, error) {
	type Alias StockOpname

	alias := &struct {
		ID           string `json:"id"`
		LocationID   string `json:"location_id"`
		CreatedByID  string `json:"created_by_id"`
		ApprovedByID string `json:"approved_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.LocationID when m.Area not nill
	// and the ID is setted
	if m.Location != nil && m.Location.ID != int64(0) {
		alias.LocationID = common.Encrypt(m.Location.ID)
	} else {
		alias.Location = nil
	}

	// Encrypt alias.CreatedByID when m.CreatedBy not nill
	// and the ID is setted
	if m.CreatedBy != nil && m.CreatedBy.ID != int64(0) {
		alias.CreatedByID = common.Encrypt(m.CreatedBy.ID)
	} else {
		alias.CreatedBy = nil
	}

	// Encrypt alias.ApprovedByID when m.ApprovedBy not nill
	// and the ID is setted
	if m.ApprovedBy != nil && m.ApprovedBy.ID != int64(0) {
		alias.ApprovedByID = common.Encrypt(m.ApprovedBy.ID)
	} else {
		alias.ApprovedBy = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating StockOpname struct into stock_opname table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to stock_opname.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *StockOpname) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.Code = m.generateCode()
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting stock_opname data
// this also will truncated all data from all table
// that have relation with this stock_opname.
func (m *StockOpname) Delete() (err error) {
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
func (m *StockOpname) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}

// generateCode to generate new stock opname code
func (m *StockOpname) generateCode() string {
	codeprefix := "STO"
	if m.Type == "adjustment" {
		codeprefix = "STA"
	}

	// format code
	now := time.Now()
	prefix := fmt.Sprintf("%s%s%s", codeprefix, now.Format("06"), now.Format("01"))

	// get last code from database
	var lastCode string
	var lastIndex int
	if e := orm.NewOrm().Raw("SELECT code from stock_opname where code like ? order by code desc", prefix+"%").QueryRow(&lastCode); e == nil {
		rn := strings.Replace(lastCode, prefix, "", 1)
		lastIndex, _ = strconv.Atoi(rn)
	}

	return m.validCode(prefix, lastIndex)
}

// validCode to check new stock opname code is valid
func (m *StockOpname) validCode(prefix string, lastIndex int) (code string) {
	lastIndex = lastIndex + 1
	code = fmt.Sprintf("%s%03d", prefix, lastIndex)

	var id int
	if e := orm.NewOrm().Raw("SELECT id FROM stock_opname where code = ?", code).QueryRow(&id); e == nil {
		code = m.validCode(prefix, lastIndex)
	}

	return
}
