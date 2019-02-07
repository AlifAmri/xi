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

	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/vehicle/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/user"
)

func init() {
	orm.RegisterModel(new(DeliveryOrder))
}

// DeliveryOrder model for DeliveryOrder table.
type DeliveryOrder struct {
	ID              int64                  `orm:"column(id);auto" json:"-"`
	Vehicle         *model.IncomingVehicle `orm:"column(vehicle_id);null;rel(fk)" json:"vehicle,omitempty"`
	Partner         *model2.Partnership    `orm:"column(partner_id);null;rel(fk)" json:"partner,omitempty"`
	Code            string                 `orm:"column(code);null;size(45)" json:"code"`
	NumberContainer string                 `orm:"column(number_container);null;size(145)" json:"number_container"`
	NumberSeal      string                 `orm:"column(number_seal);null;size(145)" json:"number_seal"`
	Status          string                 `orm:"column(status);options(active,finish)" json:"status"`
	Note            string                 `orm:"column(note);null" json:"note"`
	CreatedBy       *user.User             `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	Supply          *user.User             `orm:"column(supply_id);null;rel(fk)" json:"supply_id"`
	Checker         *user.User             `orm:"column(checker_id);null;rel(fk)" json:"checker_id"`
	DocLoc          string                 `orm:"column(docloc);null" json:"docloc"`
	CreatedAt       time.Time              `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	IsActive        int8                   `orm:"column(is_active);null" json:"is_active"`

	Tire4W  float64        `orm:"-" json:"total_tire_4w"`
	Tire2W  float64        `orm:"-" json:"total_tire_2w"`
	Counter string         `orm:"-" json:"counter_print"`
	Logo    string         `orm:"-" json:"logo"`
	Items   []*Preparation `orm:"reverse(many)" json:"items,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *DeliveryOrder) MarshalJSON() ([]byte, error) {
	type Alias DeliveryOrder

	alias := &struct {
		ID          string `json:"id"`
		CreatedByID string `json:"created_by_id"`
		VehicleID   string `json:"vehicle_id"`
		PartnerID   string `json:"partner_id"`
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

	// Encrypt alias.ApprovedByID when m.ApprovedBy not nill
	// and the ID is setted
	if m.Vehicle != nil && m.Vehicle.ID != int64(0) {
		alias.VehicleID = common.Encrypt(m.Vehicle.ID)
	} else {
		alias.Vehicle = nil
	}

	// Encrypt alias.CreatedByID when m.CreatedBy not nill
	// and the ID is setted
	if m.CreatedBy != nil && m.CreatedBy.ID != int64(0) {
		alias.CreatedByID = common.Encrypt(m.CreatedBy.ID)
	} else {
		alias.CreatedBy = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating DeliveryOrder struct into DeliveryOrder table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to DeliveryOrder.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *DeliveryOrder) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting DeliveryOrder data
// this also will truncated all data from all table
// that have relation with this DeliveryOrder.
func (m *DeliveryOrder) Delete() (err error) {
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
func (m *DeliveryOrder) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}

// Finish generate code and set status to finish
func (m *DeliveryOrder) Finish() error {
	m.Code = m.generateCode()
	m.Status = "finish"
	return m.Save("code", "status")
}

// generateCode to generate new stock opname code
func (m *DeliveryOrder) generateCode() string {
	// format code
	now := time.Now()
	suffix := fmt.Sprintf("/%s/%s/%s", "BBB", romanMonth(now.Month()), now.Format("2006"))

	// get last code from database
	var lastCode string
	var lastIndex int
	if e := orm.NewOrm().Raw("SELECT code from delivery_order where code like ? order by id desc", "%"+suffix).QueryRow(&lastCode); e == nil {
		rn := strings.Replace(lastCode, suffix, "", 1)
		lastIndex, _ = strconv.Atoi(rn)
	}

	return m.validCode(suffix, lastIndex)
}

// validCode to check new stock opname code is valid
func (m *DeliveryOrder) validCode(prefix string, lastIndex int) (code string) {
	lastIndex = lastIndex + 1
	code = fmt.Sprintf("%02d%s", lastIndex, prefix)

	var id int
	if e := orm.NewOrm().Raw("SELECT id FROM delivery_order where code = ?", code).QueryRow(&id); e == nil {
		code = m.validCode(prefix, lastIndex)
	}

	return
}

func romanMonth(m time.Month) string {
	roman := map[time.Month]string{
		1:  "I",
		2:  "II",
		3:  "III",
		4:  "IV",
		5:  "V",
		6:  "VI",
		7:  "VII",
		8:  "VIII",
		9:  "IX",
		10: "X",
		11: "XI",
		12: "XII",
	}

	return roman[m]
}
