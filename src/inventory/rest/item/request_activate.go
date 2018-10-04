// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package item

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
)

type activateRequest struct {
	ID      int64             `json:"-"`
	Item    *model.Item       `json:"-"`
	Session *auth.SessionData `json:"-"`
}

func (ar *activateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ar.Item, e = validItem(ar.ID); e != nil {
		o.Failure("id.invalid", errInvalidID)
	} else {
		if ar.Item.IsActive == 1 {
			o.Failure("id.invalid", errAlreadyActived)
		}
	}

	return o
}

func (ar *activateRequest) Messages() map[string]string {
	return map[string]string{}
}

func (ar *activateRequest) Save() error {
	ar.Item.IsActive = 1

	return ar.Item.Save("is_active")
}
