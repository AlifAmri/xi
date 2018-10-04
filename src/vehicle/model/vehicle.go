// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/partnership/model"
)

func init() {
	orm.RegisterModel(new(IncomingVehicle))
}

// IncomingVehicle model for incoming_vehicles table
type IncomingVehicle struct {
	ID              int64              `orm:"column(id);auto" json:"-"`
	DocumentID      int64              `orm:"column(document_id);null" json:"document_id"`
	Purpose         string             `orm:"column(purpose);null;options(receiving,dispatching,other)" json:"purpose"`
	Status          string             `orm:"column(status);null;options(in_progress,finished,out)" json:"status"`
	VehicleType     string             `orm:"column(vehicle_type);size(45)" json:"vehicle_type"`
	VehicleNumber   string             `orm:"column(vehicle_number);size(45)" json:"vehicle_number"`
	Driver          string             `orm:"column(driver);size(145)" json:"driver"`
	Picture         string             `orm:"column(picture);null" json:"picture"`
	CargoType       string             `orm:"column(cargo_type);options(curah,pallet,mix)" json:"cargo_type"`
	Subcon          *model.Partnership `orm:"column(subcon_id);null;rel(fk)" json:"subcon,omitempty"`
	SubconTyped     string             `orm:"column(subcon_typed);size(45);null" json:"subcon_typed"`
	Destination     string             `orm:"column(destination);options(local,export,import)" json:"destination"`
	ContainerNumber string             `orm:"column(container_number);size(45);null" json:"container_number"`
	SealNumber      string             `orm:"column(seal_number);size(45);null" json:"seal_number"`
	Notes           string             `orm:"column(notes);null" json:"notes"`
	InAt            time.Time          `orm:"column(in_at);size(45);null" json:"in_at"`
	OutAt           time.Time          `orm:"column(out_at);size(45);null" json:"out_at"`

	Receiving interface{} `orm:"-" json:"receiving,omitempty"`
	Shipment  interface{} `orm:"-" json:"shipment,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *IncomingVehicle) MarshalJSON() ([]byte, error) {
	type Alias IncomingVehicle

	alias := &struct {
		ID       string `json:"id"`
		SubconID string `json:"subcon_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.SubconID when m.Subcon not nill
	// and the ID is setted
	if m.Subcon != nil && m.Subcon.ID != int64(0) {
		alias.SubconID = common.Encrypt(m.Subcon.ID)
	} else {
		alias.Subcon = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating IncomingVehicle struct into incoming_vehicles table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to incoming_vehicles.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *IncomingVehicle) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting incoming_vehicles data
// this also will truncated all data from all table
// that have relation with this incoming_vehicles.
func (m *IncomingVehicle) Delete() (err error) {
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
func (m *IncomingVehicle) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
