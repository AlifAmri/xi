// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"time"

	"git.qasico.com/cuxs/cuxs/event"
)

func init() {
	listenLoggedin()
	listenUnauthorized()
}

func listenLoggedin() {
	c := make(chan interface{})
	event.Listen("auth::login", c)

	go func() {
		for {
			data := <-c
			user := data.(UserModelInterface)
			user.LoggedIn(time.Now())
		}
	}()
}

func listenUnauthorized() {
	c := make(chan interface{})
	event.Listen("auth::denied", c)

	go func() {
		for {
			buff := <-c
			data := buff.(*UnauthorizedLog)
			Service.UnauthorizedLog(data)
		}
	}()
}
