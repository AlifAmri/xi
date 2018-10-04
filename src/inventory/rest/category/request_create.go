// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package category

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
)

type createRequest struct {
	TypeID   string `json:"type_id" valid:"required"`
	ParentID string `json:"parent_id"`
	Name     string `json:"name" valid:"required"`
	Note     string `json:"note"`

	Session  *auth.SessionData   `json:"-"`
	ItemType *model.ItemType     `json:"-"`
	Parent   *model.ItemCategory `json:"-"`
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

		if c.ParentID != "" {
			if c.Parent, e = validParent(c.ParentID, c.ItemType.ID); e != nil {
				o.Failure("parent_id.invalid", errInvalidParent)
			}
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

func (c *createRequest) Save() (u *model.ItemCategory, e error) {
	u = &model.ItemCategory{
		Name: c.Name,
		Type: c.ItemType,
		Note: c.Note,
	}

	if c.Parent != nil {
		u.ParentID = c.Parent.ID
	}

	e = u.Save()

	return
}
