// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package partner

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/partnership/model"

	"git.qasico.com/cuxs/validation"
)

type deleteRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

func (c *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if !validID(c.ID) {
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
	t := &model.Partnership{ID: c.ID}

	return t.Delete()
}
