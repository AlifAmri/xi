// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package movement

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
	"time"
)

type startRequest struct {
	ID int64 `json:"-" valid:"required"`

	StockMovement *model.StockMovement `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

func (cr *startRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.StockMovement, e = validStockMovement(cr.ID, "new"); e != nil {
		o.Failure("id.invalid", errInvalidStockMovement)
	}

	return o
}

func (cr *startRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *startRequest) Save() (e error) {
	u := cr.StockMovement
	u.MovedBy = cr.Session.User.(*user.User)
	u.StartedAt = time.Now()
	u.Status = "start"

	if e = u.Save("moved_by", "started_at", "status"); e == nil {
		if u.IsPartial != 1 {
			u.Unit.Status = "moving"
			u.Unit.Save("status")
		}
	}

	return
}
