// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package category

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
)

type updateRequest struct {
	ID       int64  `json:"-" valid:"required"`
	ParentID string `json:"parent_id"`
	TypeID   string `json:"type_id" valid:"required"`
	Name     string `json:"name"  valid:"required"`
	Note     string `json:"note"`

	Session  *auth.SessionData   `json:"-"`
	ItemType *model.ItemType     `json:"-"`
	Parent   *model.ItemCategory `json:"-"`
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

		if c.ParentID != "" {
			if c.Parent, e = validParent(c.ParentID, c.ItemType.ID); e != nil {
				o.Failure("parent_id.invalid", errInvalidParent)
			} else {
				if c.Parent.ID == c.ID {
					o.Failure("parent_id.invalid", errInvalidSelfParent)
				}
			}
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

func (c *updateRequest) Save() (t *model.ItemCategory, e error) {
	t = &model.ItemCategory{
		ID:   c.ID,
		Type: c.ItemType,
		Name: c.Name,
		Note: c.Note,
	}

	if c.Parent != nil {
		t.ParentID = c.Parent.ID
	}

	e = t.Save()

	return
}
