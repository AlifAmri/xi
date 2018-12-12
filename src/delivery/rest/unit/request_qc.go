// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	"git.qasico.com/gudang/api/src/user"
)

type qcRequest struct {
	ID int64 `json:"-" valid:"required"`

	PreparationUnit *model.PreparationUnit `json:"-"`
	Session         *auth.SessionData      `json:"-"`
}

func (cr *qcRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.PreparationUnit, e = validPreparationUnit(cr.ID); e != nil {
		o.Failure("id.invalid", errInvalidPreparationUnit)
	}

	return o
}

func (cr *qcRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *qcRequest) Save() (e error) {
	cr.PreparationUnit.QcBy = cr.Session.User.(*user.User)
	e = cr.PreparationUnit.Save("qc_by")
	return
}
