// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/auth"

	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/cuxs/cuxs/event"
	"github.com/labstack/echo"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.Signin)
	r.GET("/me", h.Me, cuxs.Authorized())
}

// Signin endpoint to handle post http method.
func (h *Handler) Signin(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r LoginRequest
	if e = ctx.Bind(&r); e == nil {
		// provide session pada response data
		ctx.Data(auth.StartSession(r.UserData))

		// trigger event user berhasil login
		go event.Call("auth::login", r.UserData)
	}

	return ctx.Serve(e)
}

// Me http handler untuk mendapatkan session data
// user yang melakukan request
func (h *Handler) Me(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var sd *auth.SessionData
	if sd, e = auth.RequestSession(ctx); e == nil {
		ctx.Data(sd)
	}

	return ctx.Serve(e)
}
