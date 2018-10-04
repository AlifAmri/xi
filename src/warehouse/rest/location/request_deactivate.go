// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package location

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/validation"
)

type deactivateRequest struct {
	ID       int64               `json:"-"`
	Location *warehouse.Location `json:"-"`
	Session  *auth.SessionData   `json:"-"`
}

func (dr *deactivateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if dr.Location, e = validLocation(dr.ID); e != nil {
		o.Failure("id.invalid", errInvalidLocation)
	} else {
		if dr.Location.IsActive == 0 {
			o.Failure("id.invalid", errAlreadyDeactived)
		}
	}

	return o
}

func (dr *deactivateRequest) Messages() map[string]string {
	return map[string]string{}
}

func (dr *deactivateRequest) Save() error {
	dr.Location.IsActive = 0

	return dr.Location.Save("is_active")
}
