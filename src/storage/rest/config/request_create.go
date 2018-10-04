// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/storage/model"
)

type createRequest struct {
	Name      string `json:"name" valid:"required"`
	Type      string `json:"type" valid:"required"`
	TypeValue string `json:"type_value"`
	IsPrimary int8   `json:"is_primary"`
	Note      string `json:"note"`
	Areas     []area `json:"areas" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

func (cr *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if cr.Name != "" && !validName(cr.Name, 0) {
		o.Failure("name.unique", errUniqueName)
	}

	if cr.Type != "default" && cr.Type != "ncp" && cr.TypeValue == "" {
		o.Failure("type_value.required", errRequiredValue)
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":       errRequiredName,
		"type.required":       errRequiredType,
		"type_value.required": errRequiredValue,
		"areas.required":      errRequiredArea,
	}
}

func (cr *createRequest) Save() (u *model.StorageGroup, e error) {
	u = &model.StorageGroup{
		Name:      cr.Name,
		Type:      cr.Type,
		TypeValue: cr.TypeValue,
		IsActive:  1,
		IsPrimary: cr.IsPrimary,
		Note:      cr.Note,
	}

	if e = u.Save(); e == nil {
		createAreas(u, cr.Areas, false)
	}

	return
}
