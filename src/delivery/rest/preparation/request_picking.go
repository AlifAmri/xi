// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"time"

	"fmt"

	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	model3 "git.qasico.com/gudang/api/src/inventory/model"
	model2 "git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
)

type pickingRequest struct {
	ID               int64   `json:"-" valid:"required"`
	UnitID           string  `json:"unit_id" valid:"required"`
	Quantity         float64 `json:"quantity" valid:"required|gte:1"`
	QuantityRequired float64 `json:"quantity_required" valid:"required|gte:1"`

	Session     *auth.SessionData  `json:"-"`
	Preparation *model.Preparation `json:"-"`
	StockUnit   *model2.StockUnit  `json:"-"`
}

func (ur *pickingRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.Preparation, e = validPreparation(ur.ID, "active"); e != nil {
		o.Failure("id.invalid", errInvalidPreparation)
	}

	if ur.UnitID != "" {
		if ur.StockUnit, e = validStockUnit(ur.UnitID); e != nil {
			o.Failure("unit_id.invalid", errInvalidStockUnit)
		} else {
			if validLocation(ur.StockUnit.Storage.Location.ID) {
				o.Failure("location.invalid", errInvalidLocation)
			}
			if validUnitStockopname(ur.StockUnit.Storage.Location.ID) {
				o.Failure("unit_id.invalid", errInvalidStockUnitOpname)
			}
		}
	}

	if ur.StockUnit != nil && ur.StockUnit.Stock < ur.Quantity {
		message := fmt.Sprintf(", stock content hanya memiliki %v item", ur.StockUnit.Stock)
		o.Failure("quantity.invalid", errInvalidQuantity+message)
	}

	if ur.QuantityRequired < ur.Quantity {
		o.Failure("quantity.invalid", errInvalidQuantityOver)
	}
	if ur.StockUnit != nil && ur.Preparation != nil {
		if validDuplicateUnitInPreparation(ur.Preparation.ID, ur.StockUnit.ID) {
			o.Failure("unit_id.invalid", errInvalidStockUnitDuplicate)
		}
	}
	return o
}

func (ur *pickingRequest) Messages() map[string]string {
	return map[string]string{
		"unit_id.required":  errRequiredStockUnit,
		"quantity.required": errRequiredQuantity,
	}
}

func (ur *pickingRequest) Save() (u *model2.StockMovement, e error) {
	var newUnit *model2.StockUnit
	var isPartial uint8

	if ur.Quantity < ur.StockUnit.Stock {
		newUnit = new(model2.StockUnit)

		newUnit.Item = ur.StockUnit.Item
		newUnit.Batch = ur.StockUnit.Batch
		newUnit.RefID = ur.StockUnit.RefID
		newUnit.IsDefect = ur.StockUnit.IsDefect
		newUnit.Status = "draft"
		newUnit.CreatedBy = ur.StockUnit.CreatedBy
		newUnit.ReceivedAt = ur.StockUnit.ReceivedAt

		newUnit.GenerateCode(ur.StockUnit.Code)
		newUnit.Save()

		isPartial = 1
	}

	u = &model2.StockMovement{
		Unit:            ur.StockUnit,
		Type:            "picking",
		Status:          "start",
		Quantity:        ur.Quantity,
		IsPartial:       isPartial,
		Origin:          ur.StockUnit.Storage.Location,
		Destination:     ur.Preparation.Location,
		NewUnit:         newUnit,
		CreatedBy:       ur.Session.User.(*user.User),
		MovedBy:         ur.Session.User.(*user.User),
		CreatedAt:       time.Now(),
		StartedAt:       time.Now(),
		RefID:           uint64(ur.Preparation.ID),
		RefCode:         ur.Preparation.DocumentCode,
		IsNotFullPallet: int8(isPartial),
		Pallet:          &model3.Item{ID: ur.StockUnit.Storage.Container.ID},
	}

	if e = u.Save(); e == nil {
		// buat preparation unit
		pu := &model.PreparationUnit{
			Preparation:     ur.Preparation,
			Unit:            ur.StockUnit,
			Quantity:        ur.Quantity,
			LocationPicking: ur.StockUnit.Storage.Location,
			CreatedBy:       ur.Session.User.(*user.User),
			CreatedAt:       time.Now(),
		}
		if newUnit != nil {
			pu.Unit = newUnit
		}

		pu.Save()
	}

	return
}
