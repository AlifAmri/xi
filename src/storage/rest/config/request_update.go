// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/storage/model"
)

type updateRequest struct {
	ID        int64  `json:"-" valid:"required"`
	Name      string `json:"name" valid:"required"`
	Type      string `json:"type" valid:"required"`
	TypeValue string `json:"type_value" valid:"required"`
	IsPrimary int8   `json:"is_primary"`
	Note      string `json:"note"`
	Areas     []area `json:"areas" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if !validName(ur.Name, ur.ID) {
		o.Failure("name.unique", errUniqueName)
	}

	if !validID(ur.ID) {
		o.Failure("id.invalid", errInvalidStorageGroup)
	}

	if ur.Type != "default" && ur.Type != "ncp" && ur.TypeValue == "" {
		o.Failure("type_value.required", errRequiredValue)
	}

	return o
}

func (ur *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":       errRequiredName,
		"type.required":       errRequiredType,
		"type_value.required": errRequiredValue,
		"areas.required":      errRequiredArea,
	}
}

func (ur *updateRequest) Save() (u *model.StorageGroup, e error) {
	u = &model.StorageGroup{
		ID:        ur.ID,
		Name:      ur.Name,
		Type:      ur.Type,
		TypeValue: ur.TypeValue,
		IsPrimary: ur.IsPrimary,
		Note:      ur.Note,
	}

	fields := common.Fields(u, "is_active")
	if e = u.Save(fields...); e == nil {
		createAreas(u, ur.Areas, true)
	}

	return
}
