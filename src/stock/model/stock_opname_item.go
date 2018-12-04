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
	orm.RegisterModel(new(StockOpnameItem))
}

// StockOpnameItem model for stock_opname_unit table.
type StockOpnameItem struct {
	ID             int64        `orm:"column(id);auto" json:"-"`
	StockOpname    *StockOpname `orm:"column(stock_opname_id);rel(fk)" json:"stock_opname,omitempty"`
	Item           *model.Item  `orm:"column(item_id);rel(fk)" json:"item,omitempty"`
	Unit           *StockUnit   `orm:"column(unit_id);null;rel(fk)" json:"unit,omitempty"`
	UnitQuantity   float64      `orm:"column(unit_quantity);digits(12);decimals(2)" json:"unit_quantity"`
	ActualQuantity float64      `orm:"column(actual_quantity);digits(12);decimals(2)" json:"actual_quantity"`
	IsNewUnit      int8         `orm:"column(is_new_unit)" json:"is_new_unit"`
	IsDefect       int8         `orm:"column(is_defect)" json:"is_defect"`
	IsVoid         int8         `orm:"column(is_void)" json:"is_void"`
	Container      *model.Item  `orm:"column(container_id);rel(fk);null" json:"container,omitempty"`
	ContrainerNum  int8         `orm:"column(container_num)" json:"container_num"`
	Note           string       `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *StockOpnameItem) MarshalJSON() ([]byte, error) {
	type Alias StockOpnameItem

	alias := &struct {
		ID            string `json:"id"`
		CommittedByID string `json:"committed_by_id"`
		StockOpnameID string `json:"stock_opname_id"`
		ContainerID   string `json:"container_id"`
		ItemID        string `json:"item_id"`
		UnitID        string `json:"unit_id"`
		CreatedByID   string `json:"created_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.UnitID when m.Unit not nill
	// and the ID is setted
	if m.Unit != nil && m.Unit.ID != int64(0) {
		alias.UnitID = common.Encrypt(m.Unit.ID)
	} else {
		alias.Unit = nil
	}

	// Encrypt alias.StockOpnameID when m.StockOpname not nill
	// and the ID is setted
	if m.StockOpname != nil && m.StockOpname.ID != int64(0) {
		alias.StockOpnameID = common.Encrypt(m.StockOpname.ID)
	} else {
		alias.StockOpname = nil
	}

	// Encrypt alias.ItemID when m.Item not nill
	// and the ID is setted
	if m.Item != nil && m.Item.ID != int64(0) {
		alias.ItemID = common.Encrypt(m.Item.ID)
	} else {
		alias.Item = nil
	}

	if m.Container != nil && m.Container.ID != int64(0) {
		alias.ContainerID = common.Encrypt(m.Container.ID)
	} else {
		alias.Container = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating StockOpnameItem struct into stock_opname_unit table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to stock_opname_unit.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *StockOpnameItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting stock_opname_unit data
// this also will truncated all data from all table
// that have relation with this stock_opname_unit.
func (m *StockOpnameItem) Delete() (err error) {
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
func (m *StockOpnameItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
