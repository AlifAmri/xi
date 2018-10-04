// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"regexp"
	"time"

	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/usergroup"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	UsergroupID     string `json:"usergroup_id"`
	Username        string `json:"username" valid:"required"`
	Password        string `json:"password" valid:"required"`
	ConfirmPassword string `json:"confirm_password" valid:"required"`
	Name            string `json:"name" valid:"required"`
	IsSuperuser     int8   `json:"is_superuser"`

	Usergroup    *usergroup.Usergroup `json:"-"`
	Session      *auth.SessionData    `json:"-"`
	PasswordHash string               `json:"-"`
}

func (cr *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.UsergroupID == "" && cr.IsSuperuser == 0 {
		o.Failure("usergroup_id.required", errRequiredUsergroup)
	}

	// valid usergroup
	if cr.UsergroupID != "" && cr.IsSuperuser == 0 {
		if cr.Usergroup, e = validUsergroup(cr.UsergroupID); e != nil {
			o.Failure("usergroup_id.invalid", errInvalidUsergroup)
		}
	}

	if cr.Username != "" {
		if res, _ := regexp.MatchString("^[a-zA-Z0-9_]*$", cr.Username); res == false {
			o.Failure("username.invalid", errInvalidUsername)
		} else {
			if !validUsername(cr.Username, 0) {
				o.Failure("username.unique", errUniqueUsername)
			}
		}
	}

	if cr.ConfirmPassword != cr.Password {
		o.Failure("confirm_password.notmatch", errMatchPassword)
	}

	if cr.PasswordHash, e = common.PasswordHasher(cr.Password); e != nil {
		o.Failure("password.invalid", errInvalidPassword)
	}

	if cr.IsSuperuser == 1 && !validSuperuser(cr.Session) {
		o.Failure("is_superuser.invalid", errInvalidSuperuser)
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
	return map[string]string{
		"usergroup_id.required":     errRequiredUsergroup,
		"username.required":         errRequiredUsername,
		"confirm_password.required": errRequiredConfirmPassword,
		"password.required":         errRequiredPassword,
		"name.required":             errRequiredName,
	}
}

func (cr *createRequest) Save() (u *user.User, e error) {
	u = &user.User{
		Usergroup:    cr.Usergroup,
		Username:     cr.Username,
		Password:     cr.PasswordHash,
		Name:         cr.Name,
		IsActive:     1,
		IsSuperuser:  cr.IsSuperuser,
		RegisteredAt: time.Now(),
	}

	if e = u.Save(); e == nil && cr.IsSuperuser != 1 {
		auth.SetPrivilege(u.ID, u.Usergroup.ID)
	}

	return
}
