// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package partner

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/partnership/model"

	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	TypeID         string `json:"type_id" valid:"required"`
	CompanyName    string `json:"company_name" valid:"required"`
	CompanyAddress string `json:"company_address"`
	CompanyPhone   string `json:"company_phone"`
	CompanyEmail   string `json:"company_email"`
	ContactPerson  string `json:"contact_person"`
	Note           string `json:"note"`

	Session     *auth.SessionData      `json:"-"`
	PartnerType *model.PartnershipType `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	// nama harus unique
	if !validName(c.CompanyName, 0) {
		o.Failure("company_name.unique", errUniqueName)
	}

	if c.TypeID != "" {
		if c.PartnerType, e = validType(c.TypeID); e != nil {
			o.Failure("type_id.invalid", errInvalidType)
		}
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"type_id.required":      errRequiredType,
		"company_name.required": errRequiredName,
	}
}

func (c *createRequest) Save() (t *model.Partnership, e error) {
	t = &model.Partnership{
		Type:           c.PartnerType,
		CompanyName:    c.CompanyName,
		CompanyAddress: c.CompanyAddress,
		CompanyPhone:   c.CompanyPhone,
		CompanyEmail:   c.CompanyEmail,
		ContactPerson:  c.ContactPerson,
		IsActive:       1,
		Note:           c.Note,
	}

	e = t.Save()

	return
}
