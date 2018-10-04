// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/gudang/api/src/auth"
	"github.com/labstack/echo"
)

// Handler collection handler for privilege.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/item/:id", h.getItem, auth.Authorized(""))
	r.GET("/batch/:id", h.getBatch, auth.Authorized(""))
	r.GET("/unit/:id", h.getUnit, auth.Authorized(""))
}

func (h *Handler) getItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		data, total, e := GetLogItem(ctx.RequestQuery(), id)
		if e == nil {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) getBatch(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		data, total, e := GetLogBatch(ctx.RequestQuery(), id)
		if e == nil {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) getUnit(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		data, total, e := GetLogUnit(ctx.RequestQuery(), id)
		if e == nil {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}
