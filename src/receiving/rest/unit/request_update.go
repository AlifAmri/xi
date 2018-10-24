// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	model2 "git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/warehouse"
	"time"

	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID          int64   `json:"-" valid:"required"`
	LocationID  string  `json:"receiving_location_id" valid:"required"`
	UnitCode    string  `json:"unit_code"`
	ItemCode    string  `json:"item_code" valid:"required"`
	BatchCode   string  `json:"batch_code" valid:"required"`
	Quantity    float64 `json:"quantity" valid:"required|gte:1|numeric"`
	IsNcp       int8    `json:"is_ncp"`
	CheckedByID string  `json:"checked_by_id" valid:"required"`

	ReceivingUnit     *model.ReceivingUnit `json:"-"`
	ReceivingLocation *warehouse.Location  `json:"-"`
	CheckedBy         *user.User           `json:"-"`
	Session           *auth.SessionData    `json:"-"`
	Unit              *model2.StockUnit    `json:"-"`
}

func (cr *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.ReceivingUnit, e = validReceivingUnit(cr.ID, "draft"); e != nil {
		o.Failure("id.invalid", errInvalidReceivingUnit)
	}

	if cr.ReceivingUnit != nil {
		cr.ReceivingUnit.Receiving.Read()
	}

	if cr.ItemCode != "" {
		if !validItemCode(cr.ItemCode) {
			o.Failure("item_code.invalid", errInvalidItemCode)
		}
	}

	if cr.LocationID != "" {
		if cr.ReceivingLocation, e = validLocation(cr.LocationID); e != nil {
			o.Failure("receiving_location_id.invalid", errInvalidLocation)
		}
	}

	if cr.BatchCode != "" {
		if cr.BatchCode, e = validBatchCode(cr.BatchCode); e != nil {
			o.Failure("batch_code.invalid", errInvalidBatchCode)
		}
	}

	if cr.CheckedByID != "" {
		if cr.CheckedBy, e = validCheckedBy(cr.CheckedByID); e != nil {
			o.Failure("checked_by_id.invalid", errInvalidCheckedBy)
		}
	}

	if cr.UnitCode != "" && cr.ReceivingUnit != nil {
		if cr.Unit, e = validUnitCode(cr.UnitCode, cr.ID, cr.ReceivingUnit.Receiving.Plan); e != nil {
			o.Failure("unit_code.invalid", errInvalidUnit)
		}
	}

	return o
}

func (cr *updateRequest) Messages() map[string]string {
	return map[string]string{
		"receiving_location_id.required": errRequiredLocation,
		"item_code.required":             errRequiredItem,
		"batch_code.required":            errRequiredBatchCode,
		"quantity.required":              errRequiredQuantity,
		"checked_by_id.required":         errRequiredCheckedBy,
	}
}

func (cr *updateRequest) Save() (u *model.ReceivingUnit, e error) {
	u = &model.ReceivingUnit{
		ID:               cr.ReceivingUnit.ID,
		Receiving:        cr.ReceivingUnit.Receiving,
		Unit:             cr.Unit,
		UnitCode:         cr.UnitCode,
		ItemCode:         cr.ItemCode,
		BatchCode:        cr.BatchCode,
		Quantity:         cr.Quantity,
		LocationReceived: cr.ReceivingLocation,
		IsNcp:            cr.IsNcp,
		CheckedBy:        cr.CheckedBy,
		CreatedBy:        cr.Session.User.(*user.User),
		CreatedAt:        time.Now(),
	}

	e = u.Save()

	return
}
