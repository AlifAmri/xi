// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"time"

	"git.qasico.com/cuxs/cuxs/event"
	model2 "git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/stock/services/stock"
	"git.qasico.com/gudang/api/src/vehicle/model"
)

func init() {
	listenIncomingVehicle()
	listenReceivingUnitFinished()
}

func listenIncomingVehicle() {
	c := make(chan interface{})
	event.Listen("vehicle::in", c)

	go func() {
		for {
			data := <-c
			iv := data.(*model.IncomingVehicle)

			if iv.Purpose == "receiving" {
				createReceiving(iv)
			}

		}
	}()
}

func listenReceivingUnitFinished() {
	c := make(chan interface{})
	event.Listen("receiving.unit::finished", c)

	go func() {
		for {
			data := <-c
			ru := data.(*model2.ReceivingUnit)

			// create stock log
			lm := &stock.LogMaker{
				Doc:       ru.Receiving,
				Item:      ru.Unit.Item,
				Batch:     ru.Unit.Batch,
				StockUnit: ru.Unit,
				Quantity:  ru.Quantity,
			}

			stock.CreateLog(lm)
		}
	}()
}

func createReceiving(iv *model.IncomingVehicle) {
	r := &model2.Receiving{
		Vehicle:   iv,
		Status:    "active",
		StartedAt: time.Now(),
		IsActive:  1,
	}

	r.Save()
}
