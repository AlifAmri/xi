// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	partner "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/user"
)

func init() {
	orm.RegisterModel(new(ReceiptPlan))
}

// ReceiptPlan model for ReceiptPlan table.
type ReceiptPlan struct {
	ID            int64                `orm:"column(id);auto" json:"-"`
	Partner       *partner.Partnership `orm:"column(partner_id);null;rel(fk)" json:"partner,omitempty"`
	Status        string               `orm:"column(status);options(draft,pending,active,finish)" json:"status"`
	DocumentCode  string               `orm:"column(document_code);size(45)" json:"document_code"`
	TotalQuantity float64              `orm:"column(total_quantity);null;digits(12);decimals(2)" json:"total_quantity"`
	Note          string               `orm:"column(note);null" json:"note"`
	ApprovedBy    *user.User           `orm:"column(approved_by);null;rel(fk)" json:"approved_by"`
	CreatedBy     *user.User           `orm:"column(created_by);rel(fk)" json:"created_by"`
	CreatedAt     time.Time            `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	ReceivedAt    time.Time            `orm:"column(received_at);type(timestamp);null" json:"received_at"`

	Items []*ReceiptPlanItem `orm:"reverse(many)" json:"items,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ReceiptPlan) MarshalJSON() ([]byte, error) {
	type Alias ReceiptPlan

	alias := &struct {
		ID           string `json:"id"`
		ApprovedByID string `json:"approved_by_id"`
		CreatedByID  string `json:"created_by_id"`
		PartnerID    string `json:"partner_id"`
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

	return json.Marshal(alias)
}

// Save inserting or updating ReceiptPlan struct into ReceiptPlan table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to ReceiptPlan.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ReceiptPlan) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting ReceiptPlan data
// this also will truncated all data from all table
// that have relation with this ReceiptPlan.
func (m *ReceiptPlan) Delete() (err error) {
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
func (m *ReceiptPlan) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
