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
	r.PUT("/:id/finish", h.finish, auth.Authorized(""))
	r.PUT("/:id/delete", h.delete, auth.Authorized(""))
	r.PUT("/:id/qc", h.qc, auth.Authorized(""))
}

func (h *Handler) finish(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var ar finishRequest
	if ar.ID, e = ctx.Decrypt("id"); e == nil {
		if ar.Session, e = auth.RequestSession(ctx); e == nil {
			if e = ctx.Bind(&ar); e == nil {
				e = ar.Save()
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) delete(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var ar deleteRequest
	if ar.ID, e = ctx.Decrypt("id"); e == nil {
		if ar.Session, e = auth.RequestSession(ctx); e == nil {
			if e = ctx.Bind(&ar); e == nil {
				e = ar.Save()
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) qc(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var ar qcRequest
	if ar.ID, e = ctx.Decrypt("id"); e == nil {
		if ar.Session, e = auth.RequestSession(ctx); e == nil {
			if e = ctx.Bind(&ar); e == nil {
				e = ar.Save()
			}
		}
	}

	return ctx.Serve(e)
}
