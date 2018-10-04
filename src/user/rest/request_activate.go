// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/validation"
)

type activateRequest struct {
	ID      int64             `json:"-"`
	User    *user.User        `json:"-"`
	Session *auth.SessionData `json:"-"`
}

func (ar *activateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ar.User, e = validUser(ar.ID); e != nil {
		o.Failure("username.invalid", errInvalidUser)
	} else {
		if ar.User.IsSuperuser == 1 && !validSuperuser(ar.Session) {
			o.Failure("id.unauthorized", errUnauthorizedUpdate)
		}

		if ar.User.IsActive == 1 {
			o.Failure("id.invalid", errAlreadyActived)
		}
	}

	return o
}

func (ar *activateRequest) Messages() map[string]string {
	return map[string]string{}
}

func (ar *activateRequest) Save() error {
	ar.User.IsActive = 1

	return ar.User.Save("is_active")
}
