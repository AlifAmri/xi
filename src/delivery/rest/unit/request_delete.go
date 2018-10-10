// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
	"time"
)

type deleteRequest struct {
	ID int64 `json:"-" valid:"required"`

	PreparationUnit *model.PreparationUnit `json:"-"`
	Session         *auth.SessionData      `json:"-"`
}

func (cr *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.PreparationUnit, e = validPreparationUnit(cr.ID); e != nil {
		o.Failure("id.invalid", errInvalidPreparationUnit)
	}

	return o
}

func (cr *deleteRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *deleteRequest) Save() (e error) {
	if e = cr.PreparationUnit.Delete(); e == nil {
		var mv *model2.StockMovement
		e = orm.NewOrm().Raw("SELECT * FROM stock_movement "+
			"where type = 'picking' and ref_id = ? and status = 'finish' "+
			"and (unit_id = ? or new_unit = ?);", cr.PreparationUnit.Preparation.ID, cr.PreparationUnit.Unit.ID, cr.PreparationUnit.Unit.ID).QueryRow(&mv)

		if mv != nil && e == nil {
			cr.PreparationUnit.Preparation.Read()

			u := &model2.StockMovement{
				Unit:        cr.PreparationUnit.Unit,
				Type:        "routine",
				Status:      "new",
				Quantity:    cr.PreparationUnit.Quantity,
				Origin:      cr.PreparationUnit.Preparation.Location,
				Destination: mv.Origin,
				CreatedBy:   cr.Session.User.(*user.User),
				CreatedAt:   time.Now(),
			}

			if mv.NewUnit != nil {
				u.MergeUnit = mv.Unit
				u.IsMerger = uint8(1)
			}

			e = u.Save()
		}
	}

	return
}
