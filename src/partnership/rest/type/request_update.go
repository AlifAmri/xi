// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ptype

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/partnership/model"

	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID   int64  `json:"-" valid:"required"`
	Name string `json:"name"  valid:"required"`
	Note string `json:"note"`

	Session *auth.SessionData `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	// nama harus unique
	if !validName(c.Name, c.ID) {
		o.Failure("name.unique", errUniqueName)
	}

	// id harus benar
	if !validID(c.ID) {
		o.Failure("id.invalid", errInvalidID)
	}

	return o
}

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required": errRequiredName,
	}
}

func (c *updateRequest) Save() (t *model.PartnershipType, e error) {
	t = &model.PartnershipType{
		ID:   c.ID,
		Name: c.Name,
		Note: c.Note,
	}

	e = t.Save()

	return
}
