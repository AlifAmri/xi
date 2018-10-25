// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package item

import (
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/warehouse"
)

type updateRequest struct {
	ID              int64       `json:"-" valid:"required"`
	GroupID         string      `json:"group_id"`
	TypeID          string      `json:"type_id" valid:"required"`
	CategoryID      string      `json:"category_id"`
	PreferredAreaID string      `json:"preferred_area_id"`
	Code            string      `json:"code" valid:"required"`
	Name            string      `json:"name" valid:"required"`
	Image           string      `json:"image"`
	BarcodeNumber   string      `json:"barcode_number"`
	Note            string      `json:"note"`
	Attributes      []attribute `json:"attributes"`

	Session       *auth.SessionData   `json:"-"`
	ItemGroup     *model.ItemGroup    `json:"-"`
	ItemType      *model.ItemType     `json:"-"`
	ItemCategory  *model.ItemCategory `json:"-"`
	PreferredArea *warehouse.Area     `json:"-"`
	BarcodeImage  string              `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	// id harus benar
	if !validID(c.ID) {
		o.Failure("id.invalid", errInvalidID)
	}

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

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"type_id.required": errRequiredType,
		"code.required":    errRequiredCode,
		"name.required":    errRequiredName,
	}
}

func (c *updateRequest) Save() (i *model.Item, e error) {
	i = &model.Item{
		ID:            c.ID,
		Group:         c.ItemGroup,
		Type:          c.ItemType,
		Category:      c.ItemCategory,
		PreferredArea: c.PreferredArea,
		Code:          c.Code,
		Name:          c.Name,
		Image:         c.Image,
		BarcodeNumber: c.BarcodeNumber,
		BarcodeImage:  c.BarcodeImage,
		Note:          c.Note,
	}

	// update semua field kecuali stock
	fields := common.Fields(i, "stock", "is_active")
	if e = i.Save(fields...); e == nil {
		createAttribute(i, c.Attributes, true)
	}

	return
}
