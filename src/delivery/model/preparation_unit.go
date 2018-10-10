// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"git.qasico.com/gudang/api/src/warehouse"
	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
)

func init() {
	orm.RegisterModel(new(PreparationUnit))
}

// PreparationUnit model for Preparation_unit table.
type PreparationUnit struct {
	ID              int64               `orm:"column(id);auto" json:"-"`
	Preparation     *Preparation        `orm:"column(preparation_id);rel(fk)" json:"Preparation,omitempty"`
	Unit            *model.StockUnit    `orm:"column(unit_id);null;rel(fk)" json:"unit,omitempty"`
	Quantity        float64             `orm:"column(quantity);digits(12);decimals(2)" json:"quantity"`
	IsActive        int8                `orm:"column(is_active)" json:"is_active"`
	LocationPicking *warehouse.Location `orm:"column(location_picking);null;rel(fk)" json:"location_picking,omitempty"`
	CheckedBy       *user.User          `orm:"column(checked_by);null;rel(fk)" json:"checked_by"`
	CreatedBy       *user.User          `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	ApprovedBy      *user.User          `orm:"column(approved_by);null;rel(fk)" json:"approved_by"`
	CreatedAt       time.Time           `orm:"column(created_at);type(timestamp)" json:"created_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PreparationUnit) MarshalJSON() ([]byte, error) {
	type Alias PreparationUnit

	alias := &struct {
		ID                string `json:"id"`
		PreparationID     string `json:"preparation_id"`
		UnitID            string `json:"unit_id"`
		CheckedByID       string `json:"checked_by_id"`
		CreatedByID       string `json:"created_by_id"`
		ApprovedByID      string `json:"approved_by_id"`
		LocationPickingID string `json:"location_picking_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.PreparationID when m.Preparation not nill
	// and the ID is setted
	if m.Preparation != nil && m.Preparation.ID != int64(0) {
		alias.PreparationID = common.Encrypt(m.Preparation.ID)
	} else {
		alias.Preparation = nil
	}

	// Encrypt alias.UnitID when m.Unit not nill
	// and the ID is setted
	if m.Unit != nil && m.Unit.ID != int64(0) {
		alias.UnitID = common.Encrypt(m.Unit.ID)
	} else {
		alias.Unit = nil
	}

	// Encrypt alias.CheckedByID when m.CheckedBy not nill
	// and the ID is setted
	if m.CheckedBy != nil && m.CheckedBy.ID != int64(0) {
		alias.CheckedByID = common.Encrypt(m.CheckedBy.ID)
	} else {
		alias.CheckedBy = nil
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

	if m.LocationPicking != nil && m.LocationPicking.ID != int64(0) {
		alias.LocationPickingID = common.Encrypt(m.LocationPicking.ID)
	} else {
		alias.LocationPicking = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating PreparationUnit struct into Preparation_unit table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to Preparation_unit.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *PreparationUnit) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting Preparation_unit data
// this also will truncated all data from all table
// that have relation with this Preparation_unit.
func (m *PreparationUnit) Delete() (err error) {
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
func (m *PreparationUnit) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
