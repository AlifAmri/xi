// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"git.qasico.com/cuxs/cuxs/event"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
)

type finishRequest struct {
	ID int64 `json:"-" valid:"required"`

	DeliveryOrder *model.DeliveryOrder `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

func (cr *finishRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.DeliveryOrder, e = validDeliveryOrder(cr.ID, "active"); e != nil {
		o.Failure("id.invalid", errInvalidDeliveryOrder)
	}

	return o
}

func (cr *finishRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *finishRequest) Save() (e error) {
	if e = cr.DeliveryOrder.Finish(); e == nil {
		go event.Call("delivery::finished", cr.DeliveryOrder)
	}

	return
}
