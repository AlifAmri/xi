// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import "time"

// Service base model untuk authentication
//
// Memakai kit ini, saat initial run aplikasi
// harus menset model yang akan dijadikan base
//
// auth.Service = Auth{}
//
var Service BaseInterface

// BaseInterface adalah interface untuk beberapa fungsi
// auth yang harus nya dapat dicustom sesuai dengan aplikasi
type BaseInterface interface {
	// GetByID mendapatkan user data dengan berdasarkan user id
	GetByID(int64) (UserModelInterface, error)

	// GetByUsername mendapatkan user data dengan berdasarkan username
	// ini untuk kebutuhan saat login, untuk pengecekan password
	GetByUsername(string) (UserModelInterface, error)

	// UnauthorizedLog fungsi untuk menyimpan log unauthorize
	UnauthorizedLog(log *UnauthorizedLog)
}

// UserModelInterface struct interface untuk user data
// karena data user pada setiap project hampir berbeda-beda
// oleh karena itu kita menggunakan interface untuk data user
//
// Session user akan dicache saat berhasil login,
// dikarenakan cuxs/cache memakai encoding/gob
// maka saat initial run aplikasi harus meregisterkan type
// interface untuk user ini pada gob encoding.
//
// gob.Register(User{})
//
type UserModelInterface interface {
	// GetID untuk mendapatkan data id dari user tersebut
	GetID() int64

	// LoginDenied mendapatkan apakah user tersebut berhak login
	LoginAllowed() bool

	// PasswordMatched membandingkan password yang ada dengan yang diberikan
	PasswordMatched(string) bool

	// LoggedIn mengirimkan waktu last logged user
	LoggedIn(time.Time)

	// IsSuperUser apakah user ini super user
	IsSuperUser() bool
}
