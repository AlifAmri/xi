// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package area

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID   int64  `json:"-" valid:"required"`
	Code string `json:"code" valid:"required"`
	Type string `json:"type" valid:"required|in:storage,receiving,preparation,other"`
	Name string `json:"name" valid:"required"`
	Note string `json:"note"`

	Session *auth.SessionData `json:"-"`
	Area    *warehouse.Area   `json:"-"`
}

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.Area, e = validArea(ur.ID); e != nil {
		o.Failure("code.invalid", errInvalidArea)
	}

	if ur.Code != "" && !validCode(ur.Code, ur.ID) {
		o.Failure("code.unique", errUniqueCode)
	}

	return o
}

func (ur *updateRequest) Messages() map[string]string {
	return map[string]string{
		"type.required": errRequiredType,
		"type.in":       errInvalidType,
		"code.required": errRequiredCode,
		"name.required": errRequiredName,
	}
}

func (ur *updateRequest) Save() (u *warehouse.Area, e error) {
	ur.Area.Code = ur.Code
	ur.Area.Name = ur.Name
	ur.Area.Note = ur.Note
	ur.Area.Type = ur.Type

	e = ur.Area.Save("code", "name", "note", "type")

	return ur.Area, e
}
