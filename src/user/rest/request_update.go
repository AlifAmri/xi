// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"regexp"

	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/usergroup"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID              int64  `json:"-" valid:"required"`
	UsergroupID     string `json:"usergroup_id"`
	Username        string `json:"username" valid:"required"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password""`
	Name            string `json:"name" valid:"required"`
	IsSuperuser     int8   `json:"is_superuser"`

	User         *user.User           `json:"-"`
	Usergroup    *usergroup.Usergroup `json:"-"`
	Session      *auth.SessionData    `json:"-"`
	PasswordHash string               `json:"-"`
}

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.User, e = validUser(ur.ID); e != nil {
		o.Failure("username.invalid", errInvalidUser)
	}

	if ur.UsergroupID == "" && ur.IsSuperuser == 0 {
		o.Failure("usergroup_id.required", errRequiredUsergroup)
	}

	// valid usergroup
	if ur.UsergroupID != "" && ur.IsSuperuser == 0 {
		if ur.Usergroup, e = validUsergroup(ur.UsergroupID); e != nil {
			o.Failure("usergroup_id.invalid", errInvalidUsergroup)
		}
	}

	if ur.Username != "" {
		if res, _ := regexp.MatchString("^[a-zA-Z0-9_]*$", ur.Username); res == false {
			o.Failure("username.invalid", errInvalidUsername)
		} else {
			if ur.User != nil && !validUsername(ur.Username, ur.User.ID) {
				o.Failure("username.unique", errUniqueUsername)
			}
		}
	}

	if ur.Password != "" {
		if ur.ConfirmPassword != ur.Password {
			o.Failure("confirm_password.notmatch", errMatchPassword)
		}

		if ur.PasswordHash, e = common.PasswordHasher(ur.Password); e != nil {
			o.Failure("password.invalid", errInvalidPassword)
		}
	}

	if ur.IsSuperuser == 1 && !validSuperuser(ur.Session) {
		o.Failure("is_superuser.invalid", errInvalidSuperuser)
	}

	if ur.User != nil && ur.IsSuperuser == 1 && !validSuperuser(ur.Session) {
		o.Failure("username.unauthorized", errUnauthorizedUpdate)
	}

	return o
}

func (ur *updateRequest) Messages() map[string]string {
	return map[string]string{
		"usergroup_id.required":     errRequiredUsergroup,
		"username.required":         errRequiredUsername,
		"confirm_password.required": errRequiredConfirmPassword,
		"password.required":         errRequiredPassword,
		"name.required":             errRequiredName,
	}
}

func (ur *updateRequest) Save() (u *user.User, e error) {
	u = &user.User{
		ID:          ur.User.ID,
		Usergroup:   ur.Usergroup,
		Username:    ur.Username,
		Name:        ur.Name,
		IsSuperuser: ur.IsSuperuser,
	}

	var fields = []string{"usergroup_id", "username", "name", "is_superuser"}
	if ur.Password != "" {
		u.Password = ur.PasswordHash
		fields = append(fields, "password")
	}

	if e = u.Save(fields...); e == nil {
		// hapus semua privilege user ini
		orm.NewOrm().Raw("DELETE FROM privilege_user where user_id = ?", u.ID).Exec()

		if ur.IsSuperuser != 1 {
			auth.SetPrivilege(u.ID, u.Usergroup.ID)
		}
	}

	return
}
