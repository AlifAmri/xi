// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package movement

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/stock/model"
)

type deleteRequest struct {
	ID int64 `json:"-" valid:"required"`

	StockMovement *model.StockMovement `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

func (cr *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.StockMovement, e = validStockMovement(cr.ID, "new"); e != nil {
		o.Failure("id.invalid", errInvalidStockMovement)
	}

	return o
}

func (cr *deleteRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *deleteRequest) Save() (e error) {
	u := cr.StockMovement

	if u.NewUnit != nil {
		e = u.NewUnit.Delete()
	}

	if u.Type == "picking" {
		if u.NewUnit == nil {
			var puID int64
			o := orm.NewOrm()
			o.Raw("SELECT pu.id FROM preparation_unit pu "+
				"INNER JOIN preparation p ON p.id = pu.preparation_id "+
				"INNER JOIN stock_unit su ON su.id = pu.unit_id "+
				"WHERE p.id = ? AND su.id =? AND pu.is_active = ? AND pu.quantity = ?", u.RefID, u.Unit.ID, 0, u.Quantity).QueryRow(&puID)
			if puID != int64(0) {
				o.Raw("DELETE FROM preparation_unit WHERE id = ?", puID).Exec()
			}
		}
	}
	if err := u.Read("ID"); err == nil {
		e = u.Delete()
	}
	return
}
