// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/usergroup"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

var (
	errRequiredUsergroup       = "usergroup harus dipilih"
	errRequiredUsername        = "username harus diisi"
	errRequiredConfirmPassword = "konfirmasi ulang password"
	errRequiredPassword        = "password user harus diisi"
	errRequiredName            = "nama user harus diisi"
	errInvalidUsergroup        = "usergroup tidak valid"
	errInvalidUsername         = "username tidak boleh menggunakan spesial karakter"
	errUniqueUsername          = "username tersebut telah digunakan"
	errMatchPassword           = "password yang dimasukan tidak cocok"
	errInvalidPassword         = "password tidak valid"
	errInvalidSuperuser        = "Anda tidak mempunyai hak akses untuk menjadikan superuser"
	errInvalidUser             = "user tidak dapat ditemukan"
	errUnauthorizedUpdate      = "Anda tidak mempunyai hak akses untuk memperbaharui user ini"
	errAlreadyDeactived        = "status user sudah tidak aktif"
	errAlreadyActived          = "status user sudah aktif"
	errRequiredPrivilege       = "silakan pilih hak akses yang akan diberikan"
)

func validUsergroup(ide string) (ug *usergroup.Usergroup, e error) {
	ug = new(usergroup.Usergroup)
	if ug.ID, e = common.Decrypt(ide); e == nil {
		e = ug.Read()
	}

	return
}

func validUsername(username string, exclude int64) bool {
	var total int64
	orm.NewOrm().Raw("SELECT count(*) FROM user where username = ? and id != ?", username, exclude).QueryRow(&total)

	return total == 0
}

func validSuperuser(sd *auth.SessionData) bool {
	return sd.User.IsSuperUser()
}

func validUser(id int64) (u *user.User, e error) {
	u = &user.User{ID: id}
	e = u.Read()

	return
}
