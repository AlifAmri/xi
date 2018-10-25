// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package item

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	GroupID         string      `json:"group_id"`
	TypeID          string      `json:"type_id" valid:"required"`
	CategoryID      string      `json:"category_id"`
	PreferredAreaID string      `json:"preferred_area_id"`
	Code            string      `json:"code" valid:"required"`
	Name            string      `json:"name" valid:"required"`
	Image           string      `json:"image"`
	BarcodeNumber   string      `json:"barcode_number"`
	Note            string      `json:"note"`
	Attributes      []attribute `json:"Attributes"`

	Session       *auth.SessionData   `json:"-"`
	ItemGroup     *model.ItemGroup    `json:"-"`
	ItemType      *model.ItemType     `json:"-"`
	ItemCategory  *model.ItemCategory `json:"-"`
	PreferredArea *warehouse.Area     `json:"-"`
	BarcodeImage  string              `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if c.GroupID != "" {
		if c.ItemGroup, e = validGroup(c.GroupID); e != nil {
			o.Failure("group_id.invalid", errInvalidGroup)
		}
	}

	if c.TypeID != "" {
		if c.ItemType, e = validType(c.TypeID); e != nil {
			o.Failure("type_id.invalid", errInvalidType)
		}
	}

	if c.CategoryID != "" {
		if c.ItemCategory, e = validCategory(c.CategoryID); e != nil {
			o.Failure("category_id.invalid", errInvalidCategory)
		}
	}

	if c.PreferredAreaID != "" {
		if c.PreferredArea, e = validArea(c.PreferredAreaID); e != nil {
			o.Failure("preferred_area_id.invalid", errInvalidArea)
		}
	}

	// nama harus unique
	if c.ItemType != nil {
		if !validCode(c.Code, c.ItemType.ID, 0) {
			o.Failure("code.unique", errUniqueCode)
		}
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"type_id.required": errRequiredType,
		"code.required":    errRequiredCode,
		"name.required":    errRequiredName,
	}
}

func (c *createRequest) Save() (u *model.Item, e error) {
	u = &model.Item{
		Group:         c.ItemGroup,
		Type:          c.ItemType,
		Category:      c.ItemCategory,
		PreferredArea: c.PreferredArea,
		Code:          c.Code,
		Name:          c.Name,
		IsActive:      1,
		Image:         c.Image,
		BarcodeNumber: c.BarcodeNumber,
		BarcodeImage:  c.BarcodeImage,
		Note:          c.Note,
	}

	if e = u.Save(); e == nil && len(c.Attributes) > 0 {
		createAttribute(u, c.Attributes, false)
	}

	return
}
