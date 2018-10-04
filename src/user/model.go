// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"encoding/json"
	"fmt"
	"time"

	"git.qasico.com/gudang/api/src/usergroup"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(User))
}

// PrivilegeUser custom struct untuk detail user
type PrivilegeUser struct {
	Identity int64  `json:"-"`
	Name     string `json:"name"`
	Action   string `json:"action"`
	IsActive uint8  `json:"is_active"`
	Note     string `json:"note"`
	Grantted bool   `json:"grantted"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PrivilegeUser) MarshalJSON() ([]byte, error) {
	type Alias PrivilegeUser

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.Identity),
		Alias: (*Alias)(m),
	})
}

// User model for user table.
type User struct {
	ID           int64                `orm:"column(id);auto" json:"-"`
	Usergroup    *usergroup.Usergroup `orm:"column(usergroup_id);rel(fk);null" json:"usergroup,omitempty"`
	Username     string               `orm:"column(username);size(50)" json:"username"`
	Password     string               `orm:"column(password);size(150)" json:"-"`
	Name         string               `orm:"column(name);size(80)" json:"name"`
	IsActive     int8                 `orm:"column(is_active)" json:"is_active"`
	IsSuperuser  int8                 `orm:"column(is_superuser)" json:"is_superuser"`
	LastLoginAt  time.Time            `orm:"column(last_login_at);type(timestamp);null" json:"last_login_at"`
	RegisteredAt time.Time            `orm:"column(registered_at);type(timestamp)" json:"registered_at"`

	Privileges []*PrivilegeUser `orm:"-" json:"privileges,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *User) MarshalJSON() ([]byte, error) {
	type Alias User

	alias := &struct {
		ID          string `json:"id"`
		UsergroupID string `json:"usergroup_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.UsergroupID when m.Usergroup not nill
	// and the ID is setted
	if m.Usergroup != nil && m.Usergroup.ID != int64(0) {
		alias.UsergroupID = common.Encrypt(m.Usergroup.ID)
	} else {
		alias.Usergroup = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *User) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting user data
// this also will truncated all data from all table
// that have relation with this user.
func (m *User) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}

// Read execute select based on data struct that already
// assigned.
func (m *User) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}

// GetID untuk mendapatkan data id dari user tersebut
func (m *User) GetID() int64 {
	return m.ID
}

// LoginAllowed apakah user tersebut berhak login
func (m *User) LoginAllowed() bool {
	return m.IsActive == int8(1)
}

// PasswordMatched membandingkan password yang ada dengan yang diberikan
func (m *User) PasswordMatched(pwd string) bool {
	return common.PasswordHash(m.Password, pwd) == nil
}

// LoggedIn mengirimkan waktu last logged user
func (m *User) LoggedIn(ts time.Time) {
	fmt.Println("auth::login diterima di listener")
}

// IsSuperUser apakah user ini super user
func (m *User) IsSuperUser() bool {
	return m.IsSuperuser == int8(1)
}
