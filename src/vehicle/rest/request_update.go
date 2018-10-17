// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"git.qasico.com/gudang/api/src/vehicle/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	pmodel "git.qasico.com/gudang/api/src/partnership/model"
)

type updateRequest struct {
	ID              int64  `json:"-" valid:"required"`
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

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if !validID(ur.ID) {
		o.Failure("id.invalid", errInvalidID)
	}

	if ur.SubconID != "" {
		if ur.Subcon, e = validSubcon(ur.SubconID); e != nil {
			o.Failure("subcon_id.invalid", errSubconInvalid)
		}
	}

	return o
}

func (ur *updateRequest) Messages() map[string]string {
	return map[string]string{
		"vehicle_type.required":   errRequiredVehicleType,
		"vehicle_number.required": errRequiredVehicleNumber,
		"driver.required":         errRequiredDriver,
		"picture.required":        errRequiredPicture,
	}
}

func (ur *updateRequest) Save() (u *model.IncomingVehicle, e error) {
	u = &model.IncomingVehicle{
		ID:              ur.ID,
		VehicleType:     ur.VehicleType,
		VehicleNumber:   ur.VehicleNumber,
		Driver:          ur.Driver,
		Picture:         ur.Picture,
		CargoType:       "curah",
		Subcon:          ur.Subcon,
		SubconTyped:     ur.SubconTyped,
		Destination:     "local",
		ContainerNumber: ur.ContainerNumber,
		SealNumber:      ur.SealNumber,
		Notes:           ur.Notes,
	}

	fields := common.Fields(u, "in_at", "status", "purpose")
	e = u.Save(fields...)

	return
}
