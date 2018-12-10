// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package movement

import (
	"time"

	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/warehouse"
)

type createRequest struct {
	OriginID      string  `json:"origin_id" valid:"required"`
	UnitID        string  `json:"unit_id" valid:"required"`
	MergeUnitID   string  `json:"merge_unit_id"`
	DestinationID string  `json:"destination_id" valid:"required"`
	Quantity      float64 `json:"quantity" valid:"required"`
	Note          string  `json:"note"`

	Origin      *warehouse.Location `json:"-"`
	Destination *warehouse.Location `json:"-"`
	Unit        *model.StockUnit    `json:"-"`
	MergeUnit   *model.StockUnit    `json:"-"`
	IsPartial   uint8               `json:"-"`
	IsMerger    uint8               `json:"-"`
	Session     *auth.SessionData   `json:"-"`
}

func (cr *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.OriginID != "" {
		if cr.Origin, e = validLocation(cr.OriginID); e != nil {
			o.Failure("origin_id.invalid", errInvalidOrigin)
		}
	}

	if cr.UnitID != "" {
		if cr.Unit, e = validUnit(cr.UnitID); e != nil {
			o.Failure("unit_id.invalid", errInvalidUnit)
		}
	}

	if cr.DestinationID != "" {
		if cr.Destination, e = validLocation(cr.DestinationID); e != nil {
			o.Failure("destination_id.invalid", errInvalidDestination)
		}
	}

	if cr.OriginID == cr.DestinationID && cr.OriginID != "" {
		o.Failure("destination_id.invalid", errSameDestination)
	}

	if cr.Unit != nil {
		// cek apakah ada movement lain yang masih aktif untuk unit ini
		if !validUniqueMovement(cr.Unit) {
			o.Failure("unit_id.invalid", errUniqueMovement)
		}

		if cr.Origin != nil {
			if cr.Origin.ID != cr.Unit.Storage.Location.ID {
				o.Failure("unit_id.invalid", errInvalidUnit)
			}

			// cek apakah lokasi sedang di stock opname
			if !validLocationMovement(cr.Origin) {
				o.Failure("origin_id.invalid", errLocationStockopname)
			}
		}

		if cr.Quantity > cr.Unit.Stock {
			o.Failure("quantity.invalid", errInvalidQuantity)
		}

		if cr.Quantity < cr.Unit.Stock {
			cr.IsPartial = 1
		}

		if cr.MergeUnitID != "" && cr.Destination != nil {
			if !validLocationMovement(cr.Destination) {
				o.Failure("destination_id.invalid", errLocationStockopname)
			}

			if cr.MergeUnit, e = validMergeUnit(cr.MergeUnitID, cr.Unit, cr.Destination); e != nil {
				o.Failure("merge_unit_id.invalid", errInvalidMergeUnit)
			} else {
				if !validUniqueMovement(cr.MergeUnit) {
					o.Failure("merge_unit_id.invalid", errUniqueMovement)
				}

				cr.IsMerger = 1
			}
		}
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
	return map[string]string{
		"unit_id.required":        errRequiredUnit,
		"destination_id.required": errRequiredDestination,
		"quantity.required":       errRequiredQuantity,
	}
}

func (cr *createRequest) Save() (u *model.StockMovement, e error) {
	var newUnit *model.StockUnit

	if cr.IsPartial == uint8(1) && cr.MergeUnit == nil {
		newUnit = new(model.StockUnit)

		newUnit.Item = cr.Unit.Item
		newUnit.Batch = cr.Unit.Batch
		newUnit.RefID = cr.Unit.RefID
		newUnit.IsDefect = cr.Unit.IsDefect
		newUnit.Status = "draft"
		newUnit.CreatedBy = cr.Unit.CreatedBy
		newUnit.ReceivedAt = cr.Unit.ReceivedAt

		newUnit.GenerateCode(cr.Unit.Code)
		newUnit.Save()
	}

	u = &model.StockMovement{
		Unit:            cr.Unit,
		Type:            "routine",
		Status:          "new",
		Quantity:        cr.Quantity,
		IsPartial:       cr.IsPartial,
		IsMerger:        cr.IsMerger,
		Origin:          cr.Origin,
		Destination:     cr.Destination,
		NewUnit:         newUnit,
		MergeUnit:       cr.MergeUnit,
		Note:            cr.Note,
		CreatedBy:       cr.Session.User.(*user.User),
		CreatedAt:       time.Now(),
		IsNotFullPallet: int8(cr.IsPartial),
		Pallet:          cr.Unit.Storage.Container,
	}

	e = u.Save()

	return
}
