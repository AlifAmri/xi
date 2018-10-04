// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/cuxs/cuxs/event"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory"
	"git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/receiving/services"
	model3 "git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/warehouse"
)

type finishRequest struct {
	ID         int64  `json:"-" valid:"required"`
	LocationID string `json:"location_id" valid:"required"`

	ReceivingUnit *model.ReceivingUnit `json:"-"`
	Location      *warehouse.Location  `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

func (cr *finishRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.ReceivingUnit, e = validReceivingUnit(cr.ID, "draft"); e != nil {
		o.Failure("id.invalid", errInvalidReceivingUnit)
	}

	if cr.LocationID != "" {
		if cr.Location, e = validLocation(cr.LocationID); e != nil {
			o.Failure("location_id.invalid", errInvalidLocation)
		}
	}

	return o
}

func (cr *finishRequest) Messages() map[string]string {
	return map[string]string{
		"location_id.required": errRequiredLocation,
	}
}

func (cr *finishRequest) Save() (e error) {
	cr.ReceivingUnit.LocationMoved = cr.Location
	cr.ReceivingUnit.IsActive = 1
	cr.ReceivingUnit.ApprovedBy = cr.Session.User.(*user.User)
	cr.ReceivingUnit.Unit = createStockUnit(cr.ReceivingUnit)

	if e = cr.ReceivingUnit.Save("location_moved", "is_active", "approved_by", "unit_id"); e == nil {
		go event.Call("receiving.unit::finished", cr.ReceivingUnit)
		go services.CalculateActualFromUnit(cr.ReceivingUnit)
	}

	return
}

func createStockUnit(ru *model.ReceivingUnit) *model3.StockUnit {
	i := inventory.GetItem(ru.ItemCode)
	ib := inventory.GetBatch(i.ID, ru.BatchCode)

	su := &model3.StockUnit{
		Code:       ru.UnitCode,
		Item:       i,
		Batch:      ib,
		IsDefect:   ru.IsNcp,
		Status:     "moving",
		CreatedBy:  ru.ApprovedBy,
		ReceivedAt: ru.CreatedAt,
	}

	su.GenerateCode("")
	su.Save()

	return su
}
