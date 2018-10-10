// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
)

type deleteRequest struct {
	ID int64 `json:"-" valid:"required"`

	Preparation *model.Preparation `json:"-"`
	Session     *auth.SessionData  `json:"-"`
}

func (cr *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.Preparation, e = validPreparation(cr.ID, "draft"); e != nil {
		o.Failure("id.invalid", errInvalidPreparation)
	}

	return o
}

func (cr *deleteRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *deleteRequest) Save() (e error) {
	if e = cr.Preparation.Delete(); e == nil {
		if cr.Preparation.Plan != nil {
			cr.Preparation.Plan.Status = "pending"
			cr.Preparation.Plan.Save("status")
		}
	}

	return
}
