// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/auth"

	"git.qasico.com/cuxs/validation"
)

const (
	// ErrRequiredUsername error message saat username tidak diisi
	ErrRequiredUsername = "username harus diisi"

	// ErrRequiredPassword error message saat password tidak diisi
	ErrRequiredPassword = "password harus diisi"

	// ErrInvalidCredential error message saat username atau password tidak sesuai
	ErrInvalidCredential = "data tidak valid, periksa username atau password Anda kembali"

	// ErrInvalidUser error message saat system menolak user untuk login
	ErrInvalidUser = "akun Anda tidak aktif, silakan hubungi Administrator system"
)

// LoginRequest data yang harus diberikan saat login
type LoginRequest struct {
	Username string                  `valid:"required" json:"username"`
	Password string                  `valid:"required" json:"password"`
	UserData auth.UserModelInterface `json:"-"`
}

// Validate implement validation.Requests interfaces
// beriisikan custom validasi untuk request data yang diberikan
func (r *LoginRequest) Validate() *validation.Output {
	// apabila username dan password diisi
	// maka baru kita memvalidasi kedalam database
	if r.Username != "" && r.Password != "" {
		var e error
		o := &validation.Output{Valid: true}
		// memastikan bahwa username yang di berikan valid
		if r.UserData, e = validUsername(r.Username); e != nil {
			o.Failure("username.invalid", ErrInvalidCredential)
		} else {
			if !r.UserData.LoginAllowed() {
				o.Failure("username.invalid", ErrInvalidUser)
			}

			if !r.UserData.PasswordMatched(r.Password) {
				o.Failure("username.invalid", ErrInvalidCredential)
			}
		}

		return o
	}

	return nil
}

// Messages implement validation.Requests interfaces
// return custom error messages saat validasi gagal
func (r *LoginRequest) Messages() map[string]string {
	return map[string]string{
		"username.required": ErrRequiredUsername,
		"password.required": ErrRequiredPassword,
	}
}
