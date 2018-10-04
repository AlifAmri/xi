// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(ReceivingUnit))
}

// ReceivingUnit model for receiving_unit table.
type ReceivingUnit struct {
	ID                int64               `orm:"column(id);auto" json:"-"`
	Receiving         *Receiving          `orm:"column(receiving_id);rel(fk)" json:"receiving,omitempty"`
	Unit              *model.StockUnit    `orm:"column(unit_id);null;rel(fk)" json:"unit,omitempty"`
	LocationReceived  *warehouse.Location `orm:"column(location_received);null;rel(fk)" json:"location_received"`
	LocationSuggested *warehouse.Location `orm:"column(location_suggested);null;rel(fk)" json:"location_suggested"`
	LocationMoved     *warehouse.Location `orm:"column(location_moved);null;rel(fk)" json:"location_moved"`
	UnitCode          string              `orm:"column(unit_code);size(100);null" json:"unit_code"`
	ItemCode          string              `orm:"column(item_code);size(100);null" json:"item_code"`
	BatchCode         string              `orm:"column(batch_code);size(100);null" json:"batch_code"`
	Quantity          float64             `orm:"column(quantity);digits(12);decimals(2)" json:"quantity"`
	IsNcp             int8                `orm:"column(is_ncp)" json:"is_ncp"`
	IsActive          int8                `orm:"column(is_active)" json:"is_active"`
	CheckedBy         *user.User          `orm:"column(checked_by);null;rel(fk)" json:"checked_by"`
	CreatedBy         *user.User          `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	ApprovedBy        *user.User          `orm:"column(approved_by);null;rel(fk)" json:"approved_by"`
	CreatedAt         time.Time           `orm:"column(created_at);type(timestamp)" json:"created_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ReceivingUnit) MarshalJSON() ([]byte, error) {
	type Alias ReceivingUnit

	alias := &struct {
		ID                  string `json:"id"`
		ReceivingID         string `json:"receiving_id"`
		UnitID              string `json:"unit_id"`
		LocationReceivedID  string `json:"location_received_id"`
		LocationSuggestedID string `json:"location_suggested_id"`
		LocationMovedID     string `json:"location_moved_id"`
		CheckedByID         string `json:"checked_by_id"`
		CreatedByID         string `json:"created_by_id"`
		ApprovedByID        string `json:"approved_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.ReceivingID when m.Receiving not nill
	// and the ID is setted
	if m.Receiving != nil && m.Receiving.ID != int64(0) {
		alias.ReceivingID = common.Encrypt(m.Receiving.ID)
	} else {
		alias.Receiving = nil
	}

	// Encrypt alias.UnitID when m.Unit not nill
	// and the ID is setted
	if m.Unit != nil && m.Unit.ID != int64(0) {
		alias.UnitID = common.Encrypt(m.Unit.ID)
	} else {
		alias.Unit = nil
	}

	if m.LocationReceived != nil && m.LocationReceived.ID != int64(0) {
		alias.LocationReceivedID = common.Encrypt(m.LocationReceived.ID)
	} else {
		alias.LocationReceived = nil
	}

	// Encrypt alias.LocationSuggestedID when m.LocationSuggested not nill
	// and the ID is setted
	if m.LocationSuggested != nil && m.LocationSuggested.ID != int64(0) {
		alias.LocationSuggestedID = common.Encrypt(m.LocationSuggested.ID)
	} else {
		alias.LocationSuggested = nil
	}

	// Encrypt alias.LocationMovedID when m.LocationMoved not nill
	// and the ID is setted
	if m.LocationMoved != nil && m.LocationMoved.ID != int64(0) {
		alias.LocationMovedID = common.Encrypt(m.LocationMoved.ID)
	} else {
		alias.LocationMoved = nil
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

	return json.Marshal(alias)
}

// Save inserting or updating ReceivingUnit struct into receiving_unit table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to receiving_unit.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ReceivingUnit) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting receiving_unit data
// this also will truncated all data from all table
// that have relation with this receiving_unit.
func (m *ReceivingUnit) Delete() (err error) {
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
func (m *ReceivingUnit) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
