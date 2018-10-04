// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"time"

	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/cuxs/cuxs/event"
	"git.qasico.com/cuxs/orm"
	"github.com/labstack/echo"
)

// UserPrivilege mengelompokan data-data
// privilege user untuk disimpan pada session
type UserPrivilege struct {
	Name   string `json:"name"`
	Action string `json:"action"`
	Alias  string `json:"alias"`
	Note   string `json:"note"`
}

// UnauthorizedLog data-data yang akan dibutuhkan
// pada log unauthorized
type UnauthorizedLog struct {
	Action   string
	User     UserModelInterface
	LoggedAt time.Time
}

// GetPermisions mendapatkan list dari privilege
// yang dimiliki oleh user tertentu
func GetPermisions(id int64) ([]UserPrivilege, error) {
	var results []UserPrivilege

	_, e := orm.NewOrm().Raw("SELECT p.name AS name, p.action AS action, CONCAT_WS(':', p.name, p.action) AS alias, p.note AS note "+
		"FROM privilege_user pu "+
		"INNER JOIN privilege p on p.id = pu.privilege_id "+
		"WHERE pu.user_id = ? AND p.is_active = ?", id, 1).QueryRows(&results)

	return results, e
}

// SetPrivilege Usergroup mempunyai default privilege,
// maka saat user baru didaftarkan, user tersebut akan mendapatkan
// privilege default sesuai dengan usergroup nya
func SetPrivilege(userID int64, usergroupID int64) error {
	_, e := orm.NewOrm().Raw("INSERT INTO privilege_user (privilege_id, user_id) "+
		"SELECT pu.privilege_id, ? FROM privilege_usergroup as pu WHERE pu.usergroup_id = ?;", userID, usergroupID).Exec()

	return e
}

// isGranted mencek apakah alias ada pada privileges
// yang diberikan
func isGranted(alias string, privileges []UserPrivilege) (granted bool) {
	for _, privilege := range privileges {
		if privilege.Alias == alias {
			granted = true
		}
	}

	return
}

// Authorized fungsi middleware request untuk echo
// berfungsi untuk memvalidasi apakah user berhak
// mengakses suatu aksi
func Authorized(action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return cuxs.Authorized()(func(c echo.Context) error {
			// mendapatkan user session
			sd, e := RequestSession(cuxs.NewContext(c))

			// cek apakah user session diberikan
			// permisions untuk mengakses action tersebut
			if e == nil && action != "" && isGranted(action, sd.Privileges) {
				return next(c)
			}

			// hak akses untuk super user
			// super user tidak perlu data hak akses
			if e == nil && sd.User.IsSuperUser() {
				return next(c)
			}

			// hak akses tidak ada action
			if e == nil && action == "" {
				return next(c)
			}

			// trigger event user gagal atau tidak mempunyai
			// akses terhadap suatu action
			go func() {
				data := &UnauthorizedLog{
					Action:   action,
					User:     sd.User,
					LoggedAt: time.Now(),
				}

				event.Call("auth::denied", data)
			}()

			return echo.ErrUnauthorized
		})
	}
}
