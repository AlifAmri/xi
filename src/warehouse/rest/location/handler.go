// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package location

import (
	"git.qasico.com/gudang/api/src/auth"

	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for privilege.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.Authorized(""))
	r.POST("", h.create, auth.Authorized(""))
	r.GET("/:id", h.detail, auth.Authorized(""))
	r.PUT("/:id", h.update, auth.Authorized(""))
	r.PUT("/:id/activate", h.activate, auth.Authorized(""))
	r.PUT("/:id/deactivate", h.deactivate, auth.Authorized(""))
	r.PUT("/:id/delete", h.delete, auth.Authorized(""))
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

func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var cr createRequest
	if cr.Session, e = auth.RequestSession(ctx); e == nil {
		if e = ctx.Bind(&cr); e == nil {
			ctx.ResponseData, e = cr.Save()
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

func (h *Handler) activate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var ar activateRequest
	if ar.ID, e = ctx.Decrypt("id"); e == nil {
		if ar.Session, e = auth.RequestSession(ctx); e == nil {
			if e = ctx.Bind(&ar); e == nil {
				e = ar.Save()
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) deactivate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var dr deactivateRequest
	if dr.ID, e = ctx.Decrypt("id"); e == nil {
		if dr.Session, e = auth.RequestSession(ctx); e == nil {
			if e = ctx.Bind(&dr); e == nil {
				e = dr.Save()
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) delete(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var dr deleteRequest
	if dr.ID, e = ctx.Decrypt("id"); e == nil {
		if dr.Session, e = auth.RequestSession(ctx); e == nil {
			if e = ctx.Bind(&dr); e == nil {
				e = dr.Save()
			}
		}
	}

	return ctx.Serve(e)
}
