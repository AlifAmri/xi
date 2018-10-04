// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(ReceiptPlanItem))
}

// ReceiptPlanItem model for receiving_plan table.
type ReceiptPlanItem struct {
	ID        int64        `orm:"column(id);auto" json:"-"`
	Plan      *ReceiptPlan `orm:"column(plan_id);rel(fk)" json:"plan,omitempty"`
	UnitCode  string       `orm:"column(unit_code);size(100);null" json:"unit_code"`
	ItemCode  string       `orm:"column(item_code);size(100)" json:"item_code"`
	BatchCode string       `orm:"column(batch_code);size(100);null" json:"batch_code"`
	Quantity  float64      `orm:"column(quantity);digits(12);decimals(2)" json:"quantity"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ReceiptPlanItem) MarshalJSON() ([]byte, error) {
	type Alias ReceiptPlanItem

	alias := &struct {
		ID     string `json:"id"`
		PlanID string `json:"plan_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.ReceivingID when m.Receiving not nill
	// and the ID is setted
	if m.Plan != nil && m.Plan.ID != int64(0) {
		alias.PlanID = common.Encrypt(m.Plan.ID)
	} else {
		alias.Plan = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating ReceiptPlanItem struct into receiving_plan table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to receiving_plan.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ReceiptPlanItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting receiving_plan data
// this also will truncated all data from all table
// that have relation with this receiving_plan.
func (m *ReceiptPlanItem) Delete() (err error) {
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
func (m *ReceiptPlanItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
