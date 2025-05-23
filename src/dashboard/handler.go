// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dashboard

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
}

func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	if ctx.ResponseData, e = show(); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}
