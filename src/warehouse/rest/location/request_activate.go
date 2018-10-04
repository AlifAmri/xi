// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package location

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/validation"
)

type activateRequest struct {
	ID       int64               `json:"-"`
	Location *warehouse.Location `json:"-"`
	Session  *auth.SessionData   `json:"-"`
}

func (ar *activateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ar.Location, e = validLocation(ar.ID); e != nil {
		o.Failure("id.invalid", errInvalidLocation)
	} else {
		if ar.Location.IsActive == 1 {
			o.Failure("id.invalid", errAlreadyActived)
		}
	}

	return o
}

func (ar *activateRequest) Messages() map[string]string {
	return map[string]string{}
}

func (ar *activateRequest) Save() error {
	ar.Location.IsActive = 1

	return ar.Location.Save("is_active")
}
