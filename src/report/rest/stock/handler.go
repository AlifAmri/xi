// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"time"

	"git.qasico.com/cuxs/common/now"
	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/rest/batch"
	"git.qasico.com/gudang/api/src/inventory/rest/item"
	"git.qasico.com/gudang/api/src/warehouse"
	"github.com/labstack/echo"
)

// Handler collection handler for privilege.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/unit/:area_id", h.stockUnit, auth.Authorized(""))
	r.GET("/item", h.stockItem, auth.Authorized(""))
	r.GET("/batch", h.stockBatch, auth.Authorized(""))
}

// Get rest handler untuk mendapatkan data privilege
func (h *Handler) stockUnit(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	isExport := ctx.QueryParam("export") == "1"

	var id int64
	if id, e = ctx.Decrypt("area_id"); e == nil {
		wa := &warehouse.Area{ID: id}
		if e = wa.Read(); e == nil {
			data, total, e := GetUnit(rq, wa)
			if e == nil {
				if isExport {
					var file string
					if file, e = GetUnitXls(wa, *data); e == nil {
						ctx.Files(file)
					}
				} else {
					ctx.Data(data, total)
				}
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) stockItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var backdate time.Time

	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	data, total, e := item.Get(rq)
	if e == nil && total > int64(0) {
		stockItemBackDate(backdate, *data)

		if isExport {
			var file string
			if file, e = GetStockItemXls(backdate, *data); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) stockBatch(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var backdate time.Time

	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	// sort berdasarkan item dan batch
	rq.OrderBy = []string{"item_id", "code"}

	data, total, e := batch.Get(rq)
	if e == nil && total > int64(0) {
		stockBatchBackDate(backdate, *data)

		if isExport {
			var file string
			if file, e = GetStockBatchXls(backdate, *data); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}
