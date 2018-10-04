// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package storage

import (
	"git.qasico.com/cuxs/common/log"
	"git.qasico.com/cuxs/orm"
)

// Clean menghapus semua storage yang tidak dipakai
func clean() {
	_, e := orm.NewOrm().Raw("DELETE st FROM stock_storage st LEFT JOIN stock_unit su ON su.storage_id = st.id WHERE su.id IS NULL").Exec()

	log.Error(e)
}

func Recalculate() {
	clean()

	_, e := orm.NewOrm().Raw("UPDATE warehouse_location wl " +
		"SET wl.storage_used = (SELECT count(*) FROM stock_storage ss WHERE ss.location_id = wl.id) " +
		"where wl.id > 0;").Exec()

	log.Error(e)
}
