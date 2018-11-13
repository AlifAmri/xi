// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
)

type reloadRequest struct {
	ID      int64             `json:"-"`
	User    *user.User        `json:"-"`
	Session *auth.SessionData `json:"-"`
}

func (rl *reloadRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if rl.User, e = validUser(rl.ID); e != nil {
		o.Failure("username.invalid", errInvalidUser)
	} else {
		if rl.User.IsSuperuser == 1 && !validSuperuser(rl.Session) {
			o.Failure("id.unauthorized", errUnauthorizedUpdate)
		}
	}

	return o
}

func (rl *reloadRequest) Messages() map[string]string {
	return map[string]string{}
}

func (rl *reloadRequest) Reload() error {
	// hapus semua privilege user ini
	orm.NewOrm().Raw("DELETE FROM privilege_user where user_id = ?", rl.ID).Exec()

	return auth.SetPrivilege(rl.ID, rl.User.Usergroup.ID)
}
