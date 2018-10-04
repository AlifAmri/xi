// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package movement

import (
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
		u.NewUnit.Delete()
	}

	return u.Delete()
}
