// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/usergroup"

	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID         int64              `json:"-" valid:"required"`
	Name       string             `json:"name"  valid:"required"`
	Note       string             `json:"note"`
	Privileges []privilegeRequest `json:"privileges"  valid:"required"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	// nama harus unique
	if !validName(c.Name, c.ID) {
		o.Failure("name.unique", ErrUniqueName)
	}

	// id harus benar
	if !validID(c.ID) {
		o.Failure("id.invalid", ErrInvalidID)
	}

	return o
}

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":       ErrRequiredName,
		"privileges.required": ErrRequiredPrivilege,
	}
}

func (c *updateRequest) Update() (ug *usergroup.Usergroup, e error) {
	ug = &usergroup.Usergroup{
		ID:   c.ID,
		Name: c.Name,
		Note: c.Note,
	}

	if e = ug.Save(); e == nil {
		// drop privileges nya dahulu
		removePrivilege(ug.ID)
		// buat ulang privileges nya
		createPrivilege(ug.ID, c.Privileges)
		Sync(ug.ID)
	}

	return
}
