// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	"git.qasico.com/gudang/api/src/user"
)

type checkRequest struct {
	ID int64 `json:"-" valid:"required"`

	Preparation *model.Preparation `json:"-"`
	Session     *auth.SessionData  `json:"-"`
}

func (cr *checkRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	if cr.Preparation, e = validPreparationID(cr.ID); e != nil {
		o.Failure("id.invalid", "invalid preparation id")
	}

	return o
}

func (cr *checkRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *checkRequest) Save() (e error) {
	var user = cr.Session.User.(*user.User)
	cr.Preparation.CheckoutBy = user
	e = cr.Preparation.Save("checkout_by")
	return
}
