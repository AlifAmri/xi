// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	partner "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/user"
)

func init() {
	orm.RegisterModel(new(Preparation))
}

type PreparationSuggested struct {
	Location  *warehouse.Location `json:"location"`
	Item      *model.Item         `json:"item"`
	ItemBatch *model.ItemBatch    `json:"item_batch"`
	Quantity  float64             `json:"quantity"`
	YearBatch string              `json:"year"`
}

// Preparation model for Preparation table.
type Preparation struct {
	ID                  int64                `orm:"column(id);auto" json:"-"`
	Plan                *PreparationPlan     `orm:"column(plan_id);null;rel(fk)" json:"plan,omitempty"`
	Partner             *partner.Partnership `orm:"column(partner_id);null;rel(fk)" json:"partner,omitempty"`
	Supervisor          *user.User           `orm:"column(supervisor_id);null;rel(fk)" json:"supervisor,omitempty"`
	Location            *warehouse.Location  `orm:"column(location_id);null;rel(fk)" json:"location,omitempty"`
	DeliveryOrder       *DeliveryOrder       `orm:"column(delivery_order_id);null;rel(fk)" json:"delivery_order,omitempty"`
	Code                string               `orm:"column(code);null;size(45)" json:"code"`
	Status              string               `orm:"column(status);options(draft,active,finish)" json:"status"`
	DocumentCode        string               `orm:"column(document_code);size(45)" json:"document_code"`
	DocumentFile        string               `orm:"column(document_file);null" json:"document_file"`
	TotalQuantityPlan   float64              `orm:"column(total_quantity_plan);null;digits(12);decimals(2)" json:"total_quantity_plan"`
	TotalQuantityActual float64              `orm:"column(total_quantity_actual);null;digits(12);decimals(2)" json:"total_quantity_actual"`
	Note                string               `orm:"column(note);null" json:"note"`
	ApprovedBy          *user.User           `orm:"column(approved_by);null;rel(fk)" json:"approved_by"`
	CreatedBy           *user.User           `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	CreatedAt           time.Time            `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	StartedAt           time.Time            `orm:"column(started_at);type(timestamp);null" json:"started_at"`
	FinishedAt          time.Time            `orm:"column(finished_at);type(timestamp);null" json:"finished_at"`
	ShippedAt           time.Time            `orm:"column(shipped_at);type(timestamp);null" json:"shipped_at"`

	Documents []*PreparationDocument  `orm:"reverse(many)" json:"documents,omitempty"`
	Actuals   []*PreparationActual    `orm:"reverse(many)" json:"actuals,omitempty"`
	Units     []*PreparationUnit      `orm:"reverse(many)" json:"units,omitempty"`
	Pickings  []*PreparationSuggested `orm:"-" json:"pickings,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Preparation) MarshalJSON() ([]byte, error) {
	type Alias Preparation

	alias := &struct {
		ID              string `json:"id"`
		ApprovedByID    string `json:"approved_by_id"`
		CreatedByID     string `json:"created_by_id"`
		LocationID      string `json:"location_id"`
		PartnerID       string `json:"partner_id"`
		SupervisorID    string `json:"supervisor_id"`
		DeliveryOrderID string `json:"delivery_order_id"`
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

	if m.Location != nil && m.Location.ID != int64(0) {
		alias.LocationID = common.Encrypt(m.Location.ID)
	} else {
		alias.Location = nil
	}

	if m.DeliveryOrder != nil && m.DeliveryOrder.ID != int64(0) {
		alias.DeliveryOrderID = common.Encrypt(m.DeliveryOrder.ID)
	} else {
		alias.DeliveryOrder = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating Preparation struct into Preparation table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to Preparation.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Preparation) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting Preparation data
// this also will truncated all data from all table
// that have relation with this Preparation.
func (m *Preparation) Delete() (err error) {
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
func (m *Preparation) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
