// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/storage/model"
)

type deactivateRequest struct {
	ID           int64               `json:"-"`
	StorageGroup *model.StorageGroup `json:"-"`
	Session      *auth.SessionData   `json:"-"`
}

func (dr *deactivateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if dr.StorageGroup, e = validStorageGroup(dr.ID); e != nil {
		o.Failure("username.invalid", errInvalidStorageGroup)
	} else {
		if dr.StorageGroup.IsActive == 0 {
			o.Failure("id.invalid", errAlreadyDeactived)
		}
	}

	return o
}

func (dr *deactivateRequest) Messages() map[string]string {
	return map[string]string{}
}

func (dr *deactivateRequest) Save() error {
	dr.StorageGroup.IsActive = 0

	return dr.StorageGroup.Save("is_active")
}
