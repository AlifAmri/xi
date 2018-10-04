// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"fmt"
	"strings"

	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
)

type privilege struct {
	ID string `json:"id"`
}

type privilegeRequest struct {
	ID         int64       `json:"-"`
	Privileges []privilege `json:"privileges" valid:"required"`

	User *user.User `json:"-"`
}

func (pr *privilegeRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if pr.User, e = validUser(pr.ID); e != nil {
		o.Failure("id.invalid", errInvalidUser)
	} else if pr.User.IsSuperuser == 1 {
		o.Failure("id.unauthorized", errUnauthorizedUpdate)
	}

	return o
}

func (pr *privilegeRequest) Messages() map[string]string {
	return map[string]string{
		"privileges.required": errRequiredPrivilege,
	}
}

func (pr *privilegeRequest) Save() error {
	// hapus semua privilege user ini
	orm.NewOrm().Raw("DELETE FROM privilege_user where user_id = ?", pr.ID).Exec()

	var vals []string
	for _, p := range pr.Privileges {
		if pid, e := common.Decrypt(p.ID); e == nil {
			vals = append(vals, fmt.Sprintf("(%d,%d)", pid, pr.ID))
		}
	}

	_, e := orm.NewOrm().Raw("INSERT INTO privilege_user (privilege_id, user_id) VALUES " + strings.Join(vals, ",")).Exec()

	return e
}
