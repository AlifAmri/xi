// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package partner

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/partnership/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID             int64  `json:"-" valid:"required"`
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

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	// nama harus unique
	if !validName(c.CompanyName, c.ID) {
		o.Failure("name.unique", errUniqueName)
	}

	if c.TypeID != "" {
		if c.PartnerType, e = validType(c.TypeID); e != nil {
			o.Failure("type_id.invalid", errInvalidType)
		}
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

func (c *updateRequest) Save() (t *model.Partnership, e error) {
	t = &model.Partnership{
		ID:             c.ID,
		Type:           c.PartnerType,
		CompanyName:    c.CompanyName,
		CompanyAddress: c.CompanyAddress,
		CompanyPhone:   c.CompanyPhone,
		CompanyEmail:   c.CompanyEmail,
		ContactPerson:  c.ContactPerson,
		Note:           c.Note,
	}

	fields := common.Fields(t, "is_active")
	e = t.Save(fields...)

	return
}
