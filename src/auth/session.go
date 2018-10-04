// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"errors"
	"fmt"

	"git.qasico.com/cuxs/cache"
	"git.qasico.com/cuxs/cuxs"
	"github.com/dgrijalva/jwt-go"
)

// SessionData informasi-informasi mengenai user
// yang sedang login
type SessionData struct {
	Token      string             `json:"token"`
	User       UserModelInterface `json:"user"`
	Privileges []UserPrivilege    `json:"privileges"`
}

// StartSession membuat session dengan menggunakan jwt
// dan dapat menggunakan token yang lama, untuk menghindari
// pergantian token saat session.
//
// Session data akan dicache sehingga tidak perlu query ulang
// saat request selanjutnya.
func StartSession(user UserModelInterface, token ...string) (sd *SessionData) {
	if user != nil {
		sd = new(SessionData)

		// buat token baru atau menggunakan yang sebelumnya
		if len(token) == 0 {
			sd.Token = cuxs.JwtToken("id", user.GetID())
		} else {
			sd.Token = token[0]
		}

		sd.Privileges, _ = GetPermisions(user.GetID())
		sd.User = user

		cache.Set(fmt.Sprintf("session_%d", user.GetID()), sd, cache.DefaultExpiryTime)
	}

	return
}

// RequestSession mendapatkan session dari request user
// yang sedang login, akan mencoba mengambil dari cache terlebih dahulu,
// dan apabila cache tersebut telah expired, kita membuat ulang
// sessionnya dengan menggunakan token yang ada.
func RequestSession(ctx *cuxs.Context) (sd *SessionData, e error) {
	if u := ctx.Get("user"); u != nil {
		c := u.(*jwt.Token).Claims.(jwt.MapClaims)
		uid := int64(c["id"].(float64))
		if e := cache.Get(fmt.Sprintf("session_%d", uid), &sd); e != nil {
			fmt.Println(e, "==============")
			if u, e := Service.GetByID(uid); e == nil {
				token := ctx.Get("user").(*jwt.Token).Raw
				sd = StartSession(u, token)
			}
		}

		return sd, nil
	}

	return nil, errors.New("session expired")
}
