// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"fmt"
	"strings"

	"git.qasico.com/gudang/api/src/auth"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	AppName string `json:"app_name" valid:"required"`
	Company string `json:"company"  valid:"required"`
	Address string `json:"address"`
	Logo    string `json:"logo"`

	Session *auth.SessionData `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"app_name.required": errRequiredName,
		"company.required":  errRequiredCompany,
	}
}

func (c *updateRequest) Save() (v map[string]string, e error) {
	v = make(map[string]string, 4)

	v["app_name"] = c.AppName
	v["company"] = c.Company
	v["address"] = c.Address
	v["logo"] = c.Logo

	o := orm.NewOrm()
	o.Raw("TRUNCATE app_config;").Exec()

	var vals []string
	for k, v := range v {
		vals = append(vals, fmt.Sprintf("('%s','%s')", k, v))
	}

	_, e = o.Raw("INSERT INTO app_config (attribute, value) VALUES " + strings.Join(vals, ",")).Exec()

	return
}
