// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package putaway

import (
	"strconv"
	"strings"

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

	// wl storage_used harus ditambahkan dengan jumlah movement yang
	// aktif ke area tersebut
	locationMoving(wls)

	return apriory(itemCode, batchCode, wls)
}

func locationMoving(wls []*warehouse.Location) {
	var in []string
	for _, wl := range wls {
		in = append(in, strconv.Itoa(int(wl.ID)))
	}

	ins := strings.Join(in, ",")

	type x struct {
		Location int64
		Num      int
	}

	var wlm []x
	orm.NewOrm().Raw("SELECT destination_id as location, count(id) as num FROM stock_movement " +
		"where status != 'finish' and destination_id in (" + ins + ") group by destination_id;").QueryRows(&wlm)

	if len(wlm) > 0 {
		for _, xx := range wlm {
			for _, wl := range wls {
				if wl.ID == xx.Location {
					wl.StorageUsed += xx.Num
				}
			}
		}
	}
}

func apriory(itemCode string, batchCode string, wls []*warehouse.Location) (wl *warehouse.Location) {

	var pEmpty []*warehouse.Location
	var pItem []*warehouse.Location
	var pBatch []*warehouse.Location
	for _, l := range wls {
		if l.StorageUsed < l.StorageCapacity {
			if l.StorageUsed == 0 {
				pEmpty = append(pEmpty, l)
			} else {
				if batchMatch(itemCode, batchCode, l) {
					pBatch = append(pBatch, l)
					break
				}

				if itemMatch(itemCode, l) {
					pItem = append(pItem, l)
				}

				// stop looping kalau sudah ada data
				if len(pItem) > 20 {
					break
				}
			}
		}
	}

	// kalau ada yang sebatch pakai lokasi itu
	// kalau tidak ada cek lokasi kosong
	// kalau tidak ada lokasi kosong cari yang se item
	if len(pBatch) > 0 {
		wl = pBatch[0]
	} else {
		if len(pEmpty) > 0 {
			wl = pEmpty[0]
		} else if len(pItem) > 0 {
			wl = pItem[0]
		}
	}

	return
}

func batchMatch(itemCode string, batchCode string, wl *warehouse.Location) bool {
	rangeWeek := make(map[int][]int, 14)
	rangeWeek[0] = []int{1, 4}
	rangeWeek[1] = []int{5, 8}
	rangeWeek[2] = []int{9, 12}
	rangeWeek[3] = []int{13, 16}
	rangeWeek[4] = []int{17, 20}
	rangeWeek[5] = []int{21, 24}
	rangeWeek[6] = []int{25, 28}
	rangeWeek[7] = []int{29, 32}
	rangeWeek[8] = []int{33, 36}
	rangeWeek[9] = []int{37, 40}
	rangeWeek[10] = []int{41, 44}
	rangeWeek[11] = []int{45, 48}
	rangeWeek[12] = []int{49, 52}
	rangeWeek[13] = []int{53, 56}

	c, e := strconv.Atoi(batchCode[:2])
	y := batchCode[len(batchCode)-2:]

	var matched []int
	if e == nil {
		for _, x := range rangeWeek {
			if c >= x[0] && c <= x[1] {
				matched = []int{x[0], x[1]}
				break
			}
		}
	}

	if len(matched) > 0 {
		var total float64
		o := orm.NewOrm()
		o.Raw("SELECT count(*) FROM stock_unit su "+
			"inner join stock_storage ss on ss.id = su.storage_id "+
			"inner join item_batch ib on ib.id = su.batch_id "+
			"inner join item i on i.id = su.item_id "+
			"left join stock_movement sm on sm.unit_id = su.id and sm.status != 'finish' "+
			"where (SUBSTRING(ib.code, 1, 2) >= ? and SUBSTRING(ib.code, 1, 2) <= ?) "+
			"and SUBSTRING(ib.code, 3, 2) = ? and ss.location_id = ? and i.code = ? "+
			"and sm.id is null;", matched[0], matched[1], y, wl.ID, itemCode).QueryRow(&total)

		if total == 0 {
			// cek movement
			o.Raw("SELECT count(*) FROM stock_movement sm "+
				"inner join stock_unit su on su.id = sm.unit_id "+
				"inner join item_batch ib on ib.id = su.batch_id "+
				"inner join item i on i.id = su.item_id "+
				"where (SUBSTRING(ib.code, 1, 2) >= ? and SUBSTRING(ib.code, 1, 2) <= ?) "+
				"and SUBSTRING(ib.code, 3, 2) = ? and sm.destination_id = ? and sm.status != 'finish' and i.code = ?;", matched[0], matched[1], y, wl.ID, itemCode).QueryRow(&total)
		}

		return total != 0
	}

	return false
}

func itemMatch(itemCode string, wl *warehouse.Location) bool {
	var total float64
	o := orm.NewOrm()
	o.Raw("SELECT count(*) FROM stock_unit su "+
		"inner join stock_storage ss on ss.id = su.storage_id "+
		"inner join item i on i.id = su.item_id "+
		"left join stock_movement sm on sm.unit_id = su.id and sm.status != 'finish' "+
		"where i.code = ? and ss.location_id = ? "+
		"and sm.id is null;", itemCode, wl.ID).QueryRow(&total)

	if total == 0 {
		// cek movement
		o.Raw("SELECT count(*) FROM stock_movement sm "+
			"inner join stock_unit su on su.id = sm.unit_id "+
			"inner join item i on i.id = su.item_id "+
			"where i.code = ? and sm.destination_id = ? and sm.status != 'finish';", itemCode, wl.ID).QueryRow(&total)
	}

	return total != 0
}
