// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/usergroup"

	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	Name       string             `json:"name" valid:"required"`
	Note       string             `json:"note"`
	Privileges []privilegeRequest `json:"privileges" valid:"required"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	// nama harus unique
	if !validName(c.Name, 0) {
		o.Failure("name.unique", ErrUniqueName)
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":       ErrRequiredName,
		"privileges.required": ErrRequiredPrivilege,
	}
}

func (c *createRequest) Save() (ug *usergroup.Usergroup, e error) {
	ug = &usergroup.Usergroup{
		Name: c.Name,
		Note: c.Note,
	}

	if e = ug.Save(); e == nil {
		createPrivilege(ug.ID, c.Privileges)
	}

	return
}
