// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package adjustment

import (
	"time"

	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
)

type createRequest struct {
	Note  string        `json:"note"`
	Items []*opnameItem `json:"items" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

func (cr *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if len(cr.Items) > 0 {
		for i, item := range cr.Items {
			item.Validate(i, o)
		}
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
	return map[string]string{
		"items.required": errRequiredItems,
	}
}

func (cr *createRequest) Save() (u *model.StockOpname, e error) {
	u = &model.StockOpname{
		Type:      "adjustment",
		Status:    "active",
		Note:      cr.Note,
		CreatedBy: cr.Session.User.(*user.User),
		CreatedAt: time.Now(),
	}

	if e = u.Save(); e == nil {
		for _, item := range cr.Items {
			item.Save(u)
		}
	}

	return
}
