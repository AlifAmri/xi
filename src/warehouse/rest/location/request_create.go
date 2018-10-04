// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package location

import (
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	AreaID          string `json:"area_id" valid:"required"`
	Code            string `json:"code" valid:"required"`
	Name            string `json:"name" valid:"required"`
	Note            string `json:"note"`
	CordX           int    `json:"cord_x" valid:"required"`
	CordY           int    `json:"cord_y" valid:"required"`
	CordW           int    `json:"cord_w" valid:"required"`
	CordH           int    `json:"cord_h" valid:"required"`
	StorageCapacity int    `json:"storage_capacity"`

	Area    *warehouse.Area   `json:"-"`
	Session *auth.SessionData `json:"-"`
}

func (cr *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.AreaID != "" {
		if cr.Area, e = validArea(cr.AreaID); e != nil {
			o.Failure("area_id.invalid", errInvalidArea)
		}
	}

	if cr.Code != "" && cr.Area != nil && !validCode(cr.Code, cr.Area.ID, 0) {
		o.Failure("code.unique", errUniqueCode)
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
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

func (cr *createRequest) Save() (u *warehouse.Location, e error) {
	u = &warehouse.Location{
		Area:            cr.Area,
		Code:            cr.Code,
		Name:            cr.Name,
		Note:            cr.Note,
		IsActive:        1,
		CoordinateX:     cr.CordX,
		CoordinateY:     cr.CordY,
		CoordinateW:     cr.CordW,
		CoordinateH:     cr.CordH,
		StorageCapacity: cr.StorageCapacity,
	}

	e = u.Save()

	return
}
