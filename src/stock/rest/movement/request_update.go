// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package movement

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/warehouse"
)

type updateRequest struct {
	ID            int64  `json:"-" valid:"required"`
	DestinationID string `json:"destination_id" valid:"required"`
	MergeUnitID   string `json:"merge_unit_id"`

	Destination   *warehouse.Location  `json:"-"`
	MergeUnit     *model.StockUnit     `json:"-"`
	StockMovement *model.StockMovement `json:"-"`
	IsMerger      uint8                `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

func (cr *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.DestinationID != "" {
		if cr.Destination, e = validLocation(cr.DestinationID); e != nil {
			o.Failure("destination_id.invalid", errInvalidDestination)
		}
	}

	if cr.StockMovement, e = validStockMovement(cr.ID, "new"); e != nil {
		o.Failure("id.invalid", errInvalidStockMovement)
	} else {
		if cr.MergeUnitID != "" && cr.Destination != nil {
			if cr.MergeUnit, e = validMergeUnit(cr.MergeUnitID, cr.StockMovement.Unit, cr.Destination); e != nil {
				o.Failure("merge_unit_id.invalid", errInvalidMergeUnit)
			} else {
				cr.IsMerger = 1
			}
		}
	}

	return o
}

func (cr *updateRequest) Messages() map[string]string {
	return map[string]string{
		"destination_id.required": errRequiredDestination,
	}
}

func (cr *updateRequest) Save() (u *model.StockMovement, e error) {
	u = cr.StockMovement
	u.Destination = cr.Destination
	u.MergeUnit = cr.MergeUnit
	u.IsMerger = cr.IsMerger

	e = u.Save("destination_id", "merge_unit", "is_merger")

	return
}
