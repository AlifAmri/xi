// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/cuxs/cuxs/event"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	"git.qasico.com/gudang/api/src/delivery/services"
	"git.qasico.com/gudang/api/src/user"
)

type finishRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session         *auth.SessionData      `json:"-"`
	PreparationUnit *model.PreparationUnit `json:"-"`
}

func (cr *finishRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.PreparationUnit, e = validPreparationUnit(cr.ID); e != nil {
		o.Failure("id.invalid", errInvalidPreparationUnit)
	}

	return o
}

func (cr *finishRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *finishRequest) Save() (e error) {
	cr.PreparationUnit.IsActive = 1
	cr.PreparationUnit.ApprovedBy = cr.Session.User.(*user.User)
	cr.PreparationUnit.CheckedBy = cr.Session.User.(*user.User)

	if e = cr.PreparationUnit.Save("is_active", "approved_by", "checked_by"); e == nil {
		go event.Call("preparation_unit::commited", cr.PreparationUnit)
		go services.CalculateActualFromUnit(cr.PreparationUnit)
	}

	return
}
