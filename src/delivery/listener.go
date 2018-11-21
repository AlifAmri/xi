// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package delivery

import (
	"time"

	"git.qasico.com/cuxs/cuxs/event"
	model2 "git.qasico.com/gudang/api/src/delivery/model"
	"git.qasico.com/gudang/api/src/vehicle/model"
)

func init() {
	listenIncomingVehicle()
}

func listenIncomingVehicle() {
	c := make(chan interface{})
	event.Listen("vehicle::in", c)

	go func() {
		for {
			data := <-c
			iv := data.(*model.IncomingVehicle)

			if iv.Purpose == "dispatching" {
				createDeliveryOrder(iv)
			}
		}
	}()
}

func createDeliveryOrder(iv *model.IncomingVehicle) {
	r := &model2.DeliveryOrder{
		Vehicle:   iv,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	r.Save()
}
