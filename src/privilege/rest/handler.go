// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/privilege"

	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for privilege.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.Get, auth.Authorized(""))
	r.PUT("/:id/activate", h.Activate, auth.Authorized(""))
	r.PUT("/:id/deactivate", h.Deactivate, auth.Authorized(""))
}

// Get rest handler untuk mendapatkan data privilege
func (h *Handler) Get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := Get(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// Activate rest handler untuk mengaktifkan privilege
func (h *Handler) Activate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	p := new(privilege.Privilege)

	if p.ID, e = ctx.Decrypt("id"); e == nil {
		e = p.Activate()
	}

	return ctx.Serve(e)
}

// Deactivate rest handler untuk menonaktifkan privilege
func (h *Handler) Deactivate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	p := new(privilege.Privilege)

	if p.ID, e = ctx.Decrypt("id"); e == nil {
		e = p.Deactivate()
	}

	return ctx.Serve(e)
}
