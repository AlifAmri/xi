// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/cuxs/cuxs/event"
	"time"

	"git.qasico.com/gudang/api/src/vehicle/model"

	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	pmodel "git.qasico.com/gudang/api/src/partnership/model"
)

type createRequest struct {
	Purpose         string `json:"purpose" valid:"required"`
	VehicleType     string `json:"vehicle_type" valid:"required"`
	VehicleNumber   string `json:"vehicle_number" valid:"required"`
	Driver          string `json:"driver" valid:"required"`
	Picture         string `json:"picture" valid:"required"`
	CargoType       string `json:"cargo_type"`
	SubconID        string `json:"subcon_id"`
	SubconTyped     string `json:"subcon_typed"`
	Destination     string `json:"destination"`
	ContainerNumber string `json:"container_number"`
	SealNumber      string `json:"seal_number"`
	Notes           string `json:"notes"`

	Session *auth.SessionData   `json:"-"`
	Subcon  *pmodel.Partnership `json:"-"`
}

func (cr *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.SubconID != "" {
		if cr.Subcon, e = validSubcon(cr.SubconID); e != nil {
			o.Failure("subcon_id.invalid", errSubconInvalid)
		}
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
	return map[string]string{
		"purpose.required":        errRequiredPurpose,
		"vehicle_type.required":   errRequiredVehicleType,
		"vehicle_number.required": errRequiredVehicleNumber,
		"driver.required":         errRequiredDriver,
		"picture.required":        errRequiredPicture,
	}
}

func (cr *createRequest) Save() (u *model.IncomingVehicle, e error) {
	u = &model.IncomingVehicle{
		Purpose:         cr.Purpose,
		VehicleType:     cr.VehicleType,
		VehicleNumber:   cr.VehicleNumber,
		Driver:          cr.Driver,
		Picture:         cr.Picture,
		CargoType:       "curah",
		Subcon:          cr.Subcon,
		SubconTyped:     cr.SubconTyped,
		Destination:     "local",
		ContainerNumber: cr.ContainerNumber,
		SealNumber:      cr.SealNumber,
		Notes:           cr.Notes,
		Status:          "in_queue",
		InAt:            time.Now(),
	}

	if e = u.Save(); e == nil {
		go event.Call("vehicle::in", u)
	}

	return
}
