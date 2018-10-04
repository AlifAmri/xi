// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uom

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
)

type createRequest struct {
	TypeID string `json:"type_id" valid:"required"`
	Name   string `json:"name" valid:"required"`
	Note   string `json:"note"`

	Session  *auth.SessionData `json:"-"`
	ItemType *model.ItemType   `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if c.TypeID != "" {
		if c.ItemType, e = validType(c.TypeID); e != nil {
			o.Failure("type_id.invalid", errInvalidType)
		}

		if !validName(c.Name, c.ItemType.ID, 0) {
			o.Failure("name.unique", errUniqueName)
		}
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":    errRequiredName,
		"type_id.required": errRequiredType,
	}
}

func (c *createRequest) Save() (u *model.ItemUom, e error) {
	u = &model.ItemUom{
		Type: c.ItemType,
		Name: c.Name,
		Note: c.Note,
	}

	e = u.Save()

	return
}
