// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	partner "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/vehicle/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Receiving))
}

// Receiving model for receiving table.
type Receiving struct {
	ID                  int64                  `orm:"column(id);auto" json:"-"`
	Vehicle             *model.IncomingVehicle `orm:"column(vehicle_id);null;rel(fk)" json:"vehicle,omitempty"`
	Plan                *ReceiptPlan           `orm:"column(plan_id);null;rel(fk)" json:"plan,omitempty"`
	Partner             *partner.Partnership   `orm:"column(partner_id);null;rel(fk)" json:"partner,omitempty"`
	Supervisor          *user.User             `orm:"column(supervisor_id);null;rel(fk)" json:"supervisor,omitempty"`
	Code                string                 `orm:"column(code);null;size(45)" json:"code"`
	Status              string                 `orm:"column(status);options(active,finish)" json:"status"`
	DocumentCode        string                 `orm:"column(document_code);size(45)" json:"document_code"`
	DocumentFile        string                 `orm:"column(document_file);null" json:"document_file"`
	TotalQuantityPlan   float64                `orm:"column(total_quantity_plan);null;digits(12);decimals(2)" json:"total_quantity_plan"`
	TotalQuantityActual float64                `orm:"column(total_quantity_actual);null;digits(12);decimals(2)" json:"total_quantity_actual"`
	Note                string                 `orm:"column(note);null" json:"note"`
	ApprovedBy          *user.User             `orm:"column(approved_by);null;rel(fk)" json:"approved_by"`
	CreatedBy           *user.User             `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	CreatedAt           time.Time              `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	ReceivedAt          time.Time              `orm:"column(received_at);type(timestamp);null" json:"received_at"`
	StartedAt           time.Time              `orm:"column(started_at);type(timestamp);null" json:"started_at"`
	FinishedAt          time.Time              `orm:"column(finished_at);type(timestamp);null" json:"finished_at"`
	IsActive            int8                   `orm:"column(is_active);null" json:"is_active"`

	Documents []*ReceivingDocument `orm:"reverse(many)" json:"documents,omitempty"`
	Actuals   []*ReceivingActual   `orm:"reverse(many)" json:"actuals,omitempty"`
	Units     []*ReceivingUnit     `orm:"reverse(many)" json:"units,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Receiving) MarshalJSON() ([]byte, error) {
	type Alias Receiving

	alias := &struct {
		ID           string `json:"id"`
		ApprovedByID string `json:"approved_by_id"`
		CreatedByID  string `json:"created_by_id"`
		VehicleID    string `json:"vehicle_id"`
		PartnerID    string `json:"partner_id"`
		SupervisorID string `json:"supervisor_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.PartnerID when m.Partner not nill
	// and the ID is setted
	if m.Partner != nil && m.Partner.ID != int64(0) {
		alias.PartnerID = common.Encrypt(m.Partner.ID)
	} else {
		alias.Partner = nil
	}

	// Encrypt alias.SupervisorID when m.Supervisor not nill
	// and the ID is setted
	if m.Supervisor != nil && m.Supervisor.ID != int64(0) {
		alias.SupervisorID = common.Encrypt(m.Supervisor.ID)
	} else {
		alias.Supervisor = nil
	}

	// Encrypt alias.ApprovedByID when m.ApprovedBy not nill
	// and the ID is setted
	if m.ApprovedBy != nil && m.ApprovedBy.ID != int64(0) {
		alias.ApprovedByID = common.Encrypt(m.ApprovedBy.ID)
	} else {
		alias.ApprovedBy = nil
	}

	// Encrypt alias.CreatedByID when m.CreatedBy not nill
	// and the ID is setted
	if m.CreatedBy != nil && m.CreatedBy.ID != int64(0) {
		alias.CreatedByID = common.Encrypt(m.CreatedBy.ID)
	} else {
		alias.CreatedBy = nil
	}

	// Encrypt alias.VehicleID when m.Vehicle not nill
	// and the ID is setted
	if m.Vehicle != nil && m.Vehicle.ID != int64(0) {
		alias.VehicleID = common.Encrypt(m.Vehicle.ID)
	} else {
		alias.Vehicle = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating Receiving struct into receiving table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to receiving.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Receiving) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting receiving data
// this also will truncated all data from all table
// that have relation with this receiving.
func (m *Receiving) Delete() (err error) {
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
func (m *Receiving) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
