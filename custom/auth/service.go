// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/gob"
	"fmt"

	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/orm"
)

func init() {
	// setup authentication package
	auth.Service = Auth{}
	gob.Register(&user.User{})
}

// Auth base services
type Auth struct{}

// GetByID get user berdasarkan id
func (a Auth) GetByID(id int64) (u auth.UserModelInterface, e error) {
	var userModel user.User
	if e = orm.NewOrm().Raw("SELECT * FROM user where id = ?", id).QueryRow(&userModel); e == nil {
		if userModel.Usergroup != nil {
			orm.NewOrm().Raw("SELECT * FROM usergroup where id = ?", userModel.Usergroup.ID).QueryRow(&userModel.Usergroup)
		}
	}

	return &userModel, e
}

// GetByUsername get user berdasarkan username
func (a Auth) GetByUsername(username string) (auth.UserModelInterface, error) {
	var u user.User
	e := orm.NewOrm().Raw("SELECT * FROM user where username = ?", username).QueryRow(&u)

	return &u, e
}

// UnauthorizedLog untuk listener even denied
func (a Auth) UnauthorizedLog(log *auth.UnauthorizedLog) {
	fmt.Println("auth::denied diterima di listener")
}
