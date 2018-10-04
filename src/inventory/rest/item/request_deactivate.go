// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package item

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
)

type deactivateRequest struct {
	ID      int64             `json:"-"`
	Item    *model.Item       `json:"-"`
	Session *auth.SessionData `json:"-"`
}

func (dr *deactivateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if dr.Item, e = validItem(dr.ID); e != nil {
		o.Failure("id.invalid", errInvalidID)
	} else {
		if dr.Item.IsActive == 0 {
			o.Failure("id.invalid", errAlreadyDeactived)
		}
	}

	return o
}

func (dr *deactivateRequest) Messages() map[string]string {
	return map[string]string{}
}

func (dr *deactivateRequest) Save() error {
	dr.Item.IsActive = 0

	return dr.Item.Save("is_active")
}
