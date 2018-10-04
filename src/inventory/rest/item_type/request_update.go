// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package itype

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
)

type updateRequest struct {
	ID          int64  `json:"-" valid:"required"`
	Name        string `json:"name"  valid:"required"`
	IsBatch     int8   `json:"is_batch"`
	IsContainer int8   `json:"is_container"`
	Note        string `json:"note"`

	Session  *auth.SessionData `json:"-"`
	ItemType *model.ItemType   `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	// nama harus unique
	if !validName(c.Name, c.ID) {
		o.Failure("name.unique", errUniqueName)
	}

	// id harus benar
	if c.ItemType, e = validID(c.ID); e != nil {
		o.Failure("id.invalid", errInvalidID)
	} else {
		if c.ItemType.IsBatch == int8(1) && c.IsBatch == int8(0) {
			o.Failure("is_batch.invalid", errStaticBatch)
		}
		if c.ItemType.IsContainer == int8(1) && c.IsContainer == int8(0) {
			o.Failure("is_container.invalid", errStaticContainer)
		}
	}

	if c.IsContainer == 1 {
		if !validContainer(c.ID) {
			o.Failure("is_container.invalid", errInvalidContainer)
		}
	}

	return o
}

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required": errRequiredName,
	}
}

func (c *updateRequest) Save() (t *model.ItemType, e error) {
	t = &model.ItemType{
		ID:          c.ID,
		Name:        c.Name,
		IsBatch:     c.IsBatch,
		IsContainer: c.IsContainer,
		Note:        c.Note,
	}

	e = t.Save()

	return
}
