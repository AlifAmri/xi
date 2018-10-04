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
	orm.RegisterModel(new(StockMovement))
}

// StockMovement model for movement table.
type StockMovement struct {
	ID          int64               `orm:"column(id);auto" json:"-"`
	Code        string              `orm:"column(code);size(45)" json:"code"`
	Unit        *StockUnit          `orm:"column(unit_id);rel(fk)" json:"unit,omitempty"`
	Type        string              `orm:"column(type);options(routine,picking,putaway)" json:"type"`
	RefID       uint64              `orm:"column(ref_id);null" json:"-"`
	RefCode     string              `orm:"column(ref_code);null" json:"ref_code"`
	Status      string              `orm:"column(status);options(new,start,finish)" json:"status"`
	Quantity    float64             `orm:"column(quantity);digits(12);decimals(2)" json:"quantity"`
	IsPartial   uint8               `orm:"column(is_partial)" json:"is_partial"`
	IsMerger    uint8               `orm:"column(is_merger)" json:"is_merger"`
	Origin      *warehouse.Location `orm:"column(origin_id);null;rel(fk)" json:"origin,omitempty"`
	Destination *warehouse.Location `orm:"column(destination_id);null;rel(fk)" json:"destination,omitempty"`
	NewUnit     *StockUnit          `orm:"column(new_unit);null;rel(fk)" json:"new_unit"`
	MergeUnit   *StockUnit          `orm:"column(merge_unit);null;rel(fk)" json:"merge_unit"`
	Note        string              `orm:"column(note);null" json:"note"`
	CreatedBy   *user.User          `orm:"column(created_by);rel(fk)" json:"created_by"`
	MovedBy     *user.User          `orm:"column(moved_by);null;rel(fk)" json:"moved_by"`
	CreatedAt   time.Time           `orm:"column(created_at);type(timestamp)" json:"created_at"`
	StartedAt   time.Time           `orm:"column(started_at);type(timestamp);null" json:"started_at"`
	FinishedAt  time.Time           `orm:"column(finished_at);type(timestamp);null" json:"finished_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *StockMovement) MarshalJSON() ([]byte, error) {
	type Alias StockMovement

	alias := &struct {
		ID            string `json:"id"`
		RefID         string `json:"ref_id"`
		OriginID      string `json:"origin_id"`
		DestinationID string `json:"destination_id"`
		NewUnitID     string `json:"new_unit_id"`
		MergeUnitID   string `json:"merge_unit_id"`
		MovedByID     string `json:"moved_by_id"`
		CreatedByID   string `json:"created_by_id"`
		UnitID        string `json:"unit_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	if m.RefID > 0 {
		alias.RefID = common.Encrypt(m.RefID)
	}

	// Encrypt alias.CreatedByID when m.CreatedBy not nill
	// and the ID is setted
	if m.CreatedBy != nil && m.CreatedBy.ID != int64(0) {
		alias.CreatedByID = common.Encrypt(m.CreatedBy.ID)
	} else {
		alias.CreatedBy = nil
	}

	// Encrypt alias.UnitID when m.Unit not nill
	// and the ID is setted
	if m.Unit != nil && m.Unit.ID != int64(0) {
		alias.UnitID = common.Encrypt(m.Unit.ID)
	} else {
		alias.Unit = nil
	}

	// Encrypt alias.OriginID when m.Origin not nill
	// and the ID is setted
	if m.Origin != nil && m.Origin.ID != int64(0) {
		alias.OriginID = common.Encrypt(m.Origin.ID)
	} else {
		alias.Origin = nil
	}

	// Encrypt alias.DestinationID when m.Destination not nill
	// and the ID is setted
	if m.Destination != nil && m.Destination.ID != int64(0) {
		alias.DestinationID = common.Encrypt(m.Destination.ID)
	} else {
		alias.Destination = nil
	}

	// Encrypt alias.NewUnitID when m.NewUnit not nill
	// and the ID is setted
	if m.NewUnit != nil && m.NewUnit.ID != int64(0) {
		alias.NewUnitID = common.Encrypt(m.NewUnit.ID)
	} else {
		alias.NewUnit = nil
	}

	// Encrypt alias.MergeUnitID when m.MergeUnit not nill
	// and the ID is setted
	if m.MergeUnit != nil && m.MergeUnit.ID != int64(0) {
		alias.MergeUnitID = common.Encrypt(m.MergeUnit.ID)
	} else {
		alias.MergeUnit = nil
	}

	// Encrypt alias.MovedByID when m.MovedBy not nill
	// and the ID is setted
	if m.MovedBy != nil && m.MovedBy.ID != int64(0) {
		alias.MovedByID = common.Encrypt(m.MovedBy.ID)
	} else {
		alias.MovedBy = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating StockMovement struct into movement table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to movement.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *StockMovement) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.Code = m.generateCode()
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting movement data
// this also will truncated all data from all table
// that have relation with this movement.
func (m *StockMovement) Delete() (err error) {
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
func (m *StockMovement) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}

// generateCode to generate new stock opname code
func (m *StockMovement) generateCode() string {
	// format code
	codePrefix := "STM"
	now := time.Now()
	prefix := fmt.Sprintf("%s%s%s", codePrefix, now.Format("06"), now.Format("01"))
	// get last code from database
	var lastCode string
	var lastIndex int
	if e := orm.NewOrm().Raw("SELECT code from stock_movement where code like ? order by code desc", prefix+"%").QueryRow(&lastCode); e == nil {
		rn := strings.Replace(lastCode, prefix, "", 1)
		lastIndex, _ = strconv.Atoi(rn)
	}

	return m.validCode(prefix, lastIndex)
}

// validCode to check new stock opname code is valid
func (m *StockMovement) validCode(prefix string, lastIndex int) (code string) {
	lastIndex = lastIndex + 1
	code = fmt.Sprintf("%s%03d", prefix, lastIndex)

	var id int
	if e := orm.NewOrm().Raw("SELECT id FROM stock_movement where code = ?", code).QueryRow(&id); e == nil {
		code = m.validCode(prefix, lastIndex)
	}

	return
}
