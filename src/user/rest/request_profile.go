// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"regexp"

	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

type profileRequest struct {
	Username        string `json:"username" valid:"required"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Name            string `json:"name" valid:"required"`

	User         *user.User        `json:"-"`
	Session      *auth.SessionData `json:"-"`
	PasswordHash string            `json:"-"`
}

func (pr *profileRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if pr.User, e = validUser(pr.Session.User.GetID()); e != nil {
		o.Failure("username.invalid", errInvalidUser)
	}

	if pr.Username != "" {
		if res, _ := regexp.MatchString("^[a-zA-Z0-9_]*$", pr.Username); res == false {
			o.Failure("username.invalid", errInvalidUsername)
		} else {
			if pr.User != nil && !validUsername(pr.Username, pr.User.ID) {
				o.Failure("username.unique", errUniqueUsername)
			}
		}
	}

	if pr.Password != "" {
		if pr.ConfirmPassword != pr.Password {
			o.Failure("confirm_password.notmatch", errMatchPassword)
		}

		if pr.PasswordHash, e = common.PasswordHasher(pr.Password); e != nil {
			o.Failure("password.invalid", errInvalidPassword)
		}
	}

	return o
}

func (pr *profileRequest) Messages() map[string]string {
	return map[string]string{
		"username.required":         errRequiredUsername,
		"confirm_password.required": errRequiredConfirmPassword,
		"password.required":         errRequiredPassword,
		"name.required":             errRequiredName,
	}
}

func (pr *profileRequest) Save() (u *user.User, e error) {
	u = &user.User{
		ID:       pr.Session.User.GetID(),
		Username: pr.Username,
		Name:     pr.Name,
	}

	var fields = []string{"username", "name"}
	if pr.Password != "" {
		u.Password = pr.PasswordHash
		fields = append(fields, "password")
	}

	e = u.Save(fields...)

	return
}
