// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/auth"

	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for privilege.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.Get, auth.Authorized(""))
	r.POST("", h.Create, auth.Authorized(""))
	r.GET("/:id", h.Detail, auth.Authorized(""))
	r.PUT("/:id", h.Update, auth.Authorized(""))
	r.PUT("/:id/delete", h.Delete, auth.Authorized(""))
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

// Detail rest handler untuk mendapatkan data privilege
func (h *Handler) Detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = Show(id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// Create rest handler untuk mendapatkan data privilege
func (h *Handler) Create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var cr createRequest
	if e = ctx.Bind(&cr); e == nil {
		ctx.ResponseData, e = cr.Save()
	}

	return ctx.Serve(e)
}

// Update rest handler untuk mendapatkan data privilege
func (h *Handler) Update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var ur updateRequest
	if ur.ID, e = ctx.Decrypt("id"); e == nil {
		if e = ctx.Bind(&ur); e == nil {
			ctx.ResponseData, e = ur.Update()
		}
	}

	return ctx.Serve(e)
}

// Delete rest handler untuk mendapatkan data privilege
func (h *Handler) Delete(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var dr deleteRequest
	if dr.ID, e = ctx.Decrypt("id"); e == nil {
		if e = ctx.Bind(&dr); e == nil {
			e = dr.Delete()
		}
	}

	return ctx.Serve(e)
}
