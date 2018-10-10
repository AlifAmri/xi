// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"git.qasico.com/cuxs/cuxs/event"
	model3 "git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/stock/services/stock"
	"git.qasico.com/gudang/api/src/stock/services/storage"
	"time"
)

func init() {
	listenStockopnameCommited()
	listenStockopnameFinished()
	listenStockMovementFinished()
	listenPutaway()
	listenStockOut()
}

func listenStockopnameCommited() {
	c := make(chan interface{})
	event.Listen("stockopname::commited", c)

	go func() {
		for {
			data := <-c
			soi := data.(*model.StockOpnameItem)

			// create stock log
			makeLog(soi)
		}
	}()
}

func listenStockopnameFinished() {
	c := make(chan interface{})
	event.Listen("stockopname::finished", c)

	go func() {
		for {
			data := <-c
			so := data.(*model.StockOpname)

			_ = so

			// recalculate storage
			storage.Recalculate()
		}
	}()
}

func makeLog(soi *model.StockOpnameItem) {
	lm := &stock.LogMaker{
		Doc:       soi.StockOpname,
		Item:      soi.Item,
		StockUnit: soi.Unit,
		Quantity:  -1 * (soi.UnitQuantity - soi.ActualQuantity),
	}

	if soi.Unit != nil {
		lm.Batch = soi.Unit.Batch
	}

	if lm.Quantity != 0 {
		stock.CreateLog(lm)
	}
}

func listenStockMovementFinished() {
	c := make(chan interface{})
	event.Listen("stockmovement::finished", c)

	go func() {
		for {
			data := <-c
			sm := data.(*model.StockMovement)

			var lm *stock.LogMaker

			if sm.IsPartial == 1 && sm.IsMerger != 1 {
				lm = &stock.LogMaker{
					Doc:       sm,
					StockUnit: sm.NewUnit,
					Quantity:  sm.Quantity,
				}

				stock.CreateLog(lm)

				lm = &stock.LogMaker{
					Doc:       sm,
					StockUnit: sm.Unit,
					Quantity:  -1 * sm.Quantity,
				}

				stock.CreateLog(lm)
			}

			if sm.IsMerger == 1 {
				lm = &stock.LogMaker{
					Doc:       sm,
					StockUnit: sm.Unit,
					Quantity:  -1 * sm.Quantity,
				}

				stock.CreateLog(lm)

				lm = &stock.LogMaker{
					Doc:       sm,
					StockUnit: sm.MergeUnit,
					Quantity:  sm.Quantity,
				}

				stock.CreateLog(lm)
			}

			storage.Recalculate()
		}
	}()
}

func listenPutaway() {
	c := make(chan interface{})
	event.Listen("receiving.unit::finished", c)

	go func() {
		for {
			data := <-c
			ru := data.(*model2.ReceivingUnit)

			// create movement
			mv := &model.StockMovement{
				Unit:        ru.Unit,
				Type:        "putaway",
				RefID:       uint64(ru.Receiving.ID),
				RefCode:     ru.Receiving.DocumentCode,
				Status:      "new",
				Quantity:    ru.Quantity,
				Origin:      ru.LocationReceived,
				Destination: ru.LocationMoved,
				CreatedBy:   ru.ApprovedBy,
				CreatedAt:   time.Now(),
			}

			mv.Save()
		}
	}()
}

func listenStockOut() {
	c := make(chan interface{})
	event.Listen("preparation_unit::commited", c)

	go func() {
		for {
			data := <-c
			so := data.(*model3.PreparationUnit)

			so.Preparation.Read()
			so.Unit.Read()

			lm := &stock.LogMaker{
				Doc:       so.Preparation,
				Item:      so.Unit.Item,
				StockUnit: so.Unit,
				Batch:     so.Unit.Batch,
				Quantity:  -1 * so.Quantity,
			}

			if _, e := stock.CreateLog(lm); e == nil {
				so.Unit.Status = "out"
				so.Unit.Storage = nil
				so.Unit.Save("status", "storage_id")
			}

			storage.Recalculate()
		}
	}()
}
