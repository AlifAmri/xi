// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/vehicle/model"
	"time"
)

type outRequest struct {
	ID      int64                  `json:"-"`
	Session *auth.SessionData      `json:"-"`
	Vehicle *model.IncomingVehicle `json:"-"`
}

func (dr *outRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if dr.Vehicle, e = validVehicle(dr.ID); e != nil {
		o.Failure("id.invalid", errInvalidID)
	}

	return o
}

func (dr *outRequest) Messages() map[string]string {
	return map[string]string{}
}

func (dr *outRequest) Save() error {
	dr.Vehicle.OutAt = time.Now()
	dr.Vehicle.Status = "out"

	return dr.Vehicle.Save("status", "out_at")
}
