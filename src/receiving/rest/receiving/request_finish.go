// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"fmt"
	"time"

	"git.qasico.com/cuxs/cuxs/event"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/receiving/model"
)

type finishRequest struct {
	ID      int64     `json:"-" valid:"required"`
	Note    string    `json:"note"`
	Actuals []*actual `json:"items"`

	Session   *auth.SessionData `json:"-"`
	Receiving *model.Receiving  `json:"-"`
}

func (ur *finishRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.Receiving, e = validReceiving(ur.ID, "active"); e != nil {
		o.Failure("id.invalid", errInvalidReceiving)
	}

	if ur.Receiving != nil && !validFinishReceiving(ur.Receiving) {
		o.Failure("id.invalid", errInvalidReceivingProgress)
	}

	if len(ur.Actuals) > 0 {
		for i, item := range ur.Actuals {
			item.Validate(i, o)
		}
	}

	return o
}

func (ur *finishRequest) Messages() map[string]string {
	return map[string]string{}
}

func (ur *finishRequest) Save() (u *model.Receiving, e error) {
	ur.Receiving.Status = "finish"
	ur.Receiving.FinishedAt = time.Now()
	ur.Receiving.Code = fmt.Sprintf("BBB/BA-R/%06d", ur.Receiving.ID)

	if e = ur.Receiving.Save("status", "finished_at", "code"); e == nil {
		for _, item := range ur.Actuals {
			item.Save()
		}

		go event.Call("receiving::finished", ur.Receiving)
	}

	return ur.Receiving, e
}
