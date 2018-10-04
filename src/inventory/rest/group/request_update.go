// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package group

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
)

type updateRequest struct {
	ID         int64       `json:"-" valid:"required"`
	TypeID     string      `json:"type_id" valid:"required"`
	Name       string      `json:"name"  valid:"required"`
	Note       string      `json:"note"`
	Attributes []attribute `json:"Attributes"`

	Session  *auth.SessionData `json:"-"`
	ItemType *model.ItemType   `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	// id harus benar
	if !validID(c.ID) {
		o.Failure("id.invalid", errInvalidID)
	}

	if c.TypeID != "" {
		if c.ItemType, e = validType(c.TypeID); e != nil {
			o.Failure("type_id.invalid", errInvalidType)
		}

		if !validName(c.Name, c.ItemType.ID, c.ID) {
			o.Failure("name.unique", errUniqueName)
		}
	}

	return o
}

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":    errRequiredName,
		"type_id.required": errRequiredType,
	}
}

func (c *updateRequest) Save() (t *model.ItemGroup, e error) {
	t = &model.ItemGroup{
		ID:   c.ID,
		Type: c.ItemType,
		Name: c.Name,
		Note: c.Note,
	}

	if e = t.Save(); e == nil {
		createAttribute(t, c.Attributes, true)
	}

	return
}
