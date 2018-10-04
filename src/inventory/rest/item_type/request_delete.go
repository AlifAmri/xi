// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package itype

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
)

type deleteRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session  *auth.SessionData `json:"-"`
	ItemType *model.ItemType   `json:"-"`
}

func (c *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if c.ItemType, e = validID(c.ID); e != nil {
		o.Failure("id.invalid", errInvalidID)
	}

	if !validDelete(c.ID) {
		o.Failure("id.cascade", errCascadeID)
	}

	return o
}

func (c *deleteRequest) Messages() map[string]string {
	return map[string]string{}
}

func (c *deleteRequest) Save() (e error) {
	return c.ItemType.Delete()
}
