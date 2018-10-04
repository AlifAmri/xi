// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"git.qasico.com/cuxs/cuxs/event"
	"git.qasico.com/gudang/api/src/receiving/model"
)

func init() {
	listenReceivingFinished()
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
