// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
)

type publishRequest struct {
	ID int64 `json:"-" valid:"required"`

	PreparationPlan *model.PreparationPlan `json:"-"`
	Session         *auth.SessionData      `json:"-"`
}

func (cr *publishRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.PreparationPlan, e = validPreparationPlan(cr.ID, "draft"); e != nil {
		o.Failure("id.invalid", errInvalidReceiptPlan)
	}

	return o
}

func (cr *publishRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *publishRequest) Save() (e error) {
	cr.PreparationPlan.Status = "pending"

	return cr.PreparationPlan.Save("status")
}
