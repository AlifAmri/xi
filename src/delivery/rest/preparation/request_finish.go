// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"fmt"
	"git.qasico.com/cuxs/cuxs/event"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	"time"
)

type finishRequest struct {
	ID           int64     `json:"-" valid:"required"`
	Note         string    `json:"note"`
	DocumentFile string    `json:"document_file"`
	Actuals      []*actual `json:"items"`

	Session     *auth.SessionData  `json:"-"`
	Preparation *model.Preparation `json:"-"`
}

func (ur *finishRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.Preparation, e = validPreparation(ur.ID, "active"); e != nil {
		o.Failure("id.invalid", errInvalidPreparation)
	}

	if ur.Preparation != nil && !validFinishPreparation(ur.Preparation) {
		o.Failure("id.invalid", errInvalidPreparationProgress)
	}

	if len(ur.Actuals) > 0 {
		for i, item := range ur.Actuals {
			item.Validate(i, o)
		}
	}

	return o
}

func (ur *finishRequest) Messages() map[string]string {
	return map[string]string{}
}

func (ur *finishRequest) Save() (u *model.Preparation, e error) {
	ur.Preparation.DocumentFile = ur.DocumentFile
	ur.Preparation.Status = "finish"
	ur.Preparation.FinishedAt = time.Now()
	ur.Preparation.Code = fmt.Sprintf("BBB/BA-P/%06d", ur.Preparation.ID)

	if e = ur.Preparation.Save("status", "finished_at", "code"); e == nil {
		for _, item := range ur.Actuals {
			item.Save()
		}

		go event.Call("preparation::finished", ur.Preparation)
	}

	return ur.Preparation, e
}
