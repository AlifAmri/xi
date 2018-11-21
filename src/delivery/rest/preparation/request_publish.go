// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"time"

	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	"git.qasico.com/gudang/api/src/delivery/services"
	"git.qasico.com/gudang/api/src/warehouse"
)

type publishRequest struct {
	ID         int64  `json:"-" valid:"required"`
	LocationID string `json:"location_id" valid:"required"`

	Preparation *model.Preparation  `json:"-"`
	Location    *warehouse.Location `json:"-"`
	Session     *auth.SessionData   `json:"-"`
}

func (cr *publishRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.Preparation, e = validPreparation(cr.ID, "draft"); e != nil {
		o.Failure("id.invalid", errInvalidPreparation)
	}

	if cr.LocationID != "" {
		if cr.Location, e = validPreparationLocation(cr.LocationID); e != nil {
			o.Failure("location_id.invalid", errInvalidPreparationLocation)
		}
	}

	return o
}

func (cr *publishRequest) Messages() map[string]string {
	return map[string]string{
		"location_id.required": errRequiredPreparationLocation,
	}
}

func (cr *publishRequest) Save() (e error) {
	cr.Preparation.Location = cr.Location
	cr.Preparation.Status = "active"
	cr.Preparation.StartedAt = time.Now()

	if e = cr.Preparation.Save("location_id", "status", "started_at"); e == nil {
		go services.CreateActual(cr.Preparation)
	}

	return
}
