// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"git.qasico.com/cuxs/cuxs/event"
	model2 "git.qasico.com/gudang/api/src/delivery/model"
	"git.qasico.com/gudang/api/src/receiving/model"
)

func init() {
	listenReceivingFinished()
	listenDeliveryFinished()
}

func listenReceivingFinished() {
	c := make(chan interface{})
	event.Listen("receiving::finished", c)

	go func() {
		for {
			data := <-c
			so := data.(*model.Receiving)

			// update vehicle
			so.Vehicle.Status = "finished"
			so.Vehicle.Save("status")
		}
	}()
}

func listenDeliveryFinished() {
	c := make(chan interface{})
	event.Listen("delivery::finished", c)

	go func() {
		for {
			data := <-c
			so := data.(*model2.DeliveryOrder)

			// update vehicle
			so.Vehicle.Status = "finished"
			so.Vehicle.Save("status")
		}
	}()
}
