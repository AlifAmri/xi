// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package putaway

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/warehouse"
)

func SuggestedPutaway(itemCode string, batchCode string, ncp int8) (wl *warehouse.Location) {
	var wls []*warehouse.Location

	if ncp == 1 {
		wls = findLocationNCP()
	} else {
		wls = findLocationConfig(itemCode, batchCode)
	}

	return apriory(itemCode, batchCode, wls)
}

func apriory(itemCode string, batchCode string, wls []*warehouse.Location) (wl *warehouse.Location) {

	for _, l := range wls {
		if l.StorageUsed == 0 {
			wl = l
			break
		} else {
			// check apakah secode dan sebatch
			// kalau secode dan sebatch pakai lokasi ini
			if isMatch(itemCode, batchCode, l) {
				wl = l
				break
			}
		}
	}

	return
}

func isMatch(itemCode string, batchCode string, wl *warehouse.Location) bool {
	var total float64

	orm.NewOrm().Raw("SELECT count(*) FROM stock_unit su "+
		"inner join stock_storage ss on ss.id = su.storage_id "+
		"inner join item i on i.id = su.item_id "+
		"where i.code = ? and ss.location_id = ?;", itemCode, wl.ID).QueryRow(&total)

	return total != 0
}
