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

type createRequest struct {
	ReceivingID string  `json:"receiving_id" valid:"required"`
	LocationID  string  `json:"receiving_location_id" valid:"required"`
	UnitCode    string  `json:"unit_code"`
	ItemCode    string  `json:"item_code" valid:"required"`
	BatchCode   string  `json:"batch_code" valid:"required"`
	Quantity    float64 `json:"quantity" valid:"required|gte:1|numeric"`
	IsNcp       int8    `json:"is_ncp"`
	CheckedByID string  `json:"checked_by_id" valid:"required"`

	Receiving         *model.Receiving    `json:"-"`
	ReceivingLocation *warehouse.Location `json:"-"`
	CheckedBy         *user.User          `json:"-"`
	Session           *auth.SessionData   `json:"-"`
	Unit              *model2.StockUnit   `json:"-"`
}

func (cr *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.ReceivingID != "" {
		if cr.Receiving, e = validReceiving(cr.ReceivingID); e != nil {
			o.Failure("receiving_id.invalid", errInvalidReceiving)
		}
	}

	if cr.BatchCode != "" {
		if cr.BatchCode, e = validBatchCode(cr.BatchCode); e != nil {
			o.Failure("batch_code.invalid", errInvalidBatchCode)
		}
	}

	if cr.LocationID != "" {
		if cr.ReceivingLocation, e = validLocation(cr.LocationID); e != nil {
			o.Failure("receiving_location_id.invalid", errInvalidLocation)
		}
	}

	if cr.CheckedByID != "" {
		if cr.CheckedBy, e = validCheckedBy(cr.CheckedByID); e != nil {
			o.Failure("checked_by_id.invalid", errInvalidCheckedBy)
		}
	}

	if cr.UnitCode != "" && cr.Receiving != nil {
		if cr.Unit, e = validUnitCode(cr.UnitCode, 0, cr.Receiving.Plan); e != nil {
			o.Failure("unit_code.invalid", errInvalidUnit)
		}
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
	return map[string]string{
		"receiving_id.required":          errRequiredReceiving,
		"receiving_location_id.required": errRequiredLocation,
		"item_code.required":             errRequiredItem,
		"batch_code.required":            errRequiredBatchCode,
		"quantity.required":              errRequiredQuantity,
		"checked_by_id.required":         errRequiredCheckedBy,
	}
}

func (cr *createRequest) Save() (u *model.ReceivingUnit, e error) {
	u = &model.ReceivingUnit{
		Unit:             cr.Unit,
		Receiving:        cr.Receiving,
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
