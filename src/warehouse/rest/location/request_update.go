// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package location

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID              int64  `json:"-" valid:"required"`
	AreaID          string `json:"area_id" valid:"required"`
	Code            string `json:"code" valid:"required"`
	Name            string `json:"name" valid:"required"`
	Note            string `json:"note"`
	CordX           int    `json:"cord_x" valid:"required"`
	CordY           int    `json:"cord_y" valid:"required"`
	CordW           int    `json:"cord_w" valid:"required"`
	CordH           int    `json:"cord_h" valid:"required"`
	StorageCapacity int    `json:"storage_capacity"`

	Session  *auth.SessionData   `json:"-"`
	Location *warehouse.Location `json:"-"`
	Area     *warehouse.Area     `json:"-"`
}

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.Location, e = validLocation(ur.ID); e != nil {
		o.Failure("code.invalid", errInvalidLocation)
	}

	if ur.AreaID != "" {
		if ur.Area, e = validArea(ur.AreaID); e != nil {
			o.Failure("area_id.invalid", errInvalidArea)
		}
	}

	if ur.Code != "" && ur.Area != nil && !validCode(ur.Code, ur.Area.ID, ur.Location.ID) {
		o.Failure("code.unique", errUniqueCode)
	}

	return o
}

func (ur *updateRequest) Messages() map[string]string {
	return map[string]string{
		"area_id.required": errRequiredArea,
		"code.required":    errRequiredCode,
		"name.required":    errRequiredName,
		"cord_x.required":  errRequiredCordX,
		"cord_y.required":  errRequiredCordY,
		"cord_w.required":  errRequiredCordW,
		"cord_h.required":  errRequiredCordH,
	}
}

func (ur *updateRequest) Save() (u *warehouse.Location, e error) {
	ur.Location.Code = ur.Code
	ur.Location.Name = ur.Name
	ur.Location.Note = ur.Note
	ur.Location.CoordinateX = ur.CordX
	ur.Location.CoordinateY = ur.CordY
	ur.Location.CoordinateW = ur.CordW
	ur.Location.CoordinateH = ur.CordH
	ur.Location.StorageCapacity = ur.StorageCapacity

	e = ur.Location.Save()

	return ur.Location, e
}
