// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package area

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	Code string `json:"code" valid:"required"`
	Type string `json:"type" valid:"required|in:storage,receiving,preparation,other"`
	Name string `json:"name" valid:"required"`
	Note string `json:"note"`

	Session *auth.SessionData `json:"-"`
}

func (cr *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if cr.Code != "" && !validCode(cr.Code, 0) {
		o.Failure("code.unique", errUniqueCode)
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
	return map[string]string{
		"type.required": errRequiredType,
		"type.in":       errInvalidType,
		"code.required": errRequiredCode,
		"name.required": errRequiredName,
	}
}

func (cr *createRequest) Save() (u *warehouse.Area, e error) {
	u = &warehouse.Area{
		Warehouse: &warehouse.Warehouse{ID: 1},
		Code:      cr.Code,
		Type:      cr.Type,
		Name:      cr.Name,
		Note:      cr.Note,
		IsActive:  1,
	}

	e = u.Save()

	return
}
