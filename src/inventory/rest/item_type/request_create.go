// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package itype

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
)

type createRequest struct {
	Name        string `json:"name" valid:"required"`
	IsBatch     int8   `json:"is_batch"`
	IsContainer int8   `json:"is_container"`
	Note        string `json:"note"`

	Session *auth.SessionData `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	// nama harus unique
	if !validName(c.Name, 0) {
		o.Failure("name.unique", errUniqueName)
	}

	if c.IsContainer == 1 {
		if !validContainer(0) {
			o.Failure("is_container.invalid", errInvalidContainer)
		}
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required": errRequiredName,
	}
}

func (c *createRequest) Save() (t *model.ItemType, e error) {
	t = &model.ItemType{
		Name:        c.Name,
		IsBatch:     c.IsBatch,
		IsContainer: c.IsContainer,
		Note:        c.Note,
	}

	e = t.Save()

	return
}
