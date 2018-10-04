// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"fmt"
	"strings"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

const (
	// ErrRequiredName error message saat nama tidak diisi
	ErrRequiredName = "nama harus diisi"

	// ErrRequiredPrivilege error message saat privilege tidak diisi
	ErrRequiredPrivilege = "privilege harus dipilih"

	// ErrUniqueName error message saat nama yang diberikan sudah digunakan
	ErrUniqueName = "nama tersebut sudah digunakan"

	// ErrInvalidID error message saat id yang digunakan tidak benar
	ErrInvalidID = "identitas data tidak ditemukan"

	// ErrCascadeID error message saat usergroup masih digunakan oleh user
	ErrCascadeID = "usergroup masih digunakan oleh user"
)

type privilegeRequest struct {
	ID string `json:"id" valid:"required"`
}

func createPrivilege(id int64, pr []privilegeRequest) {
	var vals []string
	for _, p := range pr {
		if pid, e := common.Decrypt(p.ID); e == nil {
			vals = append(vals, fmt.Sprintf("(%d,%d)", pid, id))
		}
	}

	orm.NewOrm().Raw("INSERT INTO privilege_usergroup (privilege_id, usergroup_id) VALUES " + strings.Join(vals, ",")).Exec()
}

func removePrivilege(id int64) {
	orm.NewOrm().Raw("DELETE FROM privilege_usergroup where usergroup_id = ?", id).Exec()
}
