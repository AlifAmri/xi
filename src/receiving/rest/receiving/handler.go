// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/auth"
	"github.com/labstack/echo"
)

// Handler collection handler for privilege.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.Authorized(""))
	r.GET("/:id", h.detail, auth.Authorized(""))
	r.PUT("/sync/:id", h.sync, auth.Authorized(""))
	r.PUT("/:id", h.update, auth.Authorized(""))
	r.PUT("/finish/:id", h.finish, auth.Authorized(""))
	r.PUT("/out/:id", h.out, auth.Authorized(""))
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

func (h *Handler) sync(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var ur syncRequest
	if ur.ID, e = ctx.Decrypt("id"); e == nil {
		if ur.Session, e = auth.RequestSession(ctx); e == nil {
			if e = ctx.Bind(&ur); e == nil {
				e = ur.Save()
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var ur updateRequest
	if ur.ID, e = ctx.Decrypt("id"); e == nil {
		if ur.Session, e = auth.RequestSession(ctx); e == nil {
			if e = ctx.Bind(&ur); e == nil {
				ctx.ResponseData, e = ur.Save()
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) finish(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var ur finishRequest
	if ur.ID, e = ctx.Decrypt("id"); e == nil {
		if ur.Session, e = auth.RequestSession(ctx); e == nil {
			if e = ctx.Bind(&ur); e == nil {
				ctx.ResponseData, e = ur.Save()
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) out(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var ID int64
	if ID, e = ctx.Decrypt("id"); e == nil {
		_, e = orm.NewOrm().Raw("UPDATE incoming_vehicle ivh INNER JOIN receiving r ON r.vehicle_id = ivh.id "+
			"SET ivh.status = 'finished' WHERE r.id = ?", ID).Exec()
		if e == nil {
			ctx.ResponseData = "success"
		}
	}

	return ctx.Serve(e)
}
