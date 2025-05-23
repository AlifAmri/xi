// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/gudang/api/src/auth"
	"github.com/labstack/echo"
)

// Handler collection handler for privilege.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.Authorized(""))
	r.GET("/:id", h.detail, auth.Authorized(""))
	r.GET("/:id/moving", h.moving, auth.Authorized(""))
	r.GET("/:id/prepare", h.prepare, auth.Authorized(""))
	r.GET("/:id/receiving", h.receiving, auth.Authorized(""))
	r.GET("/:id/opname", h.opname, auth.Authorized(""))
}

func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	data, total, e := Get(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = Show(id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) moving(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		data, total, e := historyMovement(id, ctx.RequestQuery())
		if e == nil {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) prepare(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		data, total, e := historyPreparation(id, ctx.RequestQuery())
		if e == nil {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) receiving(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		data, total, e := historyReceiving(id, ctx.RequestQuery())
		if e == nil {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) opname(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		data, total, e := historyStockopname(id, ctx.RequestQuery())
		if e == nil {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}
