// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/vehicle/model"
)

type cancelRequest struct {
	ID int64 `json:"-" valid:"required"`

	Vehicle *model.IncomingVehicle `json:"-"`
	Session *auth.SessionData      `json:"-"`
}

func (cr *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.Vehicle, e = validCancel(cr.ID, "in_queue"); e != nil {
		o.Failure("id.invalid", errCancelInvalid)
	}

	return o
}

func (cr *cancelRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *cancelRequest) Save() (e error) {
	u := cr.Vehicle
	o := orm.NewOrm()
	u.Status = "cancelled"
	if e = u.Save("status"); e == nil {
		if u.Purpose == "receiving" {
			_, e = o.Raw("update receiving set is_active = ? where vehicle_id = ?", 0, u.ID).Exec()
		}

		if u.Purpose == "dispatching" {
			_, e = o.Raw("update delivery_order set is_active = ?, status = ? where vehicle_id = ?", 0 , "cancelled", u.ID).Exec()
		}
	}

	return
}
