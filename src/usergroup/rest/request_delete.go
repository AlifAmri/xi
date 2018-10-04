// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/usergroup"

	"git.qasico.com/cuxs/validation"
)

type deleteRequest struct {
	ID int64 `json:"-" valid:"required"`
}

func (c *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if !validID(c.ID) {
		o.Failure("id.invalid", ErrInvalidID)
	}

	if !validDelete(c.ID) {
		o.Failure("id.cascade", ErrCascadeID)
	}

	return o
}

func (c *deleteRequest) Messages() map[string]string {
	return map[string]string{}
}

func (c *deleteRequest) Delete() (e error) {
	ug := &usergroup.Usergroup{ID: c.ID}

	if e = ug.Delete(); e == nil {
		removePrivilege(ug.ID)
	}

	return
}
