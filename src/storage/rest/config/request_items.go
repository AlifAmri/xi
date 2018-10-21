// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"strings"

	"git.qasico.com/gudang/api/src/storage/model"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
)

type item struct {
	AreaID         string `json:"area_id"  valid:"required"`
	LocationFromID string `json:"location_from_id" valid:"required"`
	LocationEndID  string `json:"location_end_id"  valid:"required"`

	Area         *warehouse.Area     `json:"-"`
	LocationFrom *warehouse.Location `json:"-"`
	LocationEnd  *warehouse.Location `json:"-"`
}

func (oi *item) Validate(index int, o *validation.Output, SgID int64) {
	var e error

	if oi.AreaID != "" {
		if oi.Area, e = validArea(oi.AreaID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.area_id.invalid", index), errInvalidArea)
		}

		if oi.LocationFromID != "" && oi.Area != nil {
			if oi.LocationFrom, e = validLocationFrom(oi.LocationFromID, oi.Area); e != nil {
				o.Failure(fmt.Sprintf("items.%d.location_from_id.invalid", index), errInvalidLocation)
			}
		}

		if oi.LocationEndID != "" && oi.LocationFrom != nil && oi.Area != nil {
			if oi.LocationEnd, e = validLocationEnd(oi.LocationEndID, oi.Area, oi.LocationFrom); e != nil {
				o.Failure(fmt.Sprintf("items.%d.location_end_id.invalid", index), errInvalidLocation)
			}
		}

		if oi.LocationFrom != nil {
			if !validLocation(oi.LocationFrom, SgID) {
				o.Failure(fmt.Sprintf("items.%d.location_from_id.unique", index), errUniqueLocation)
			}
		}

		if oi.LocationEnd != nil {
			if !validLocation(oi.LocationEnd, SgID) {
				o.Failure(fmt.Sprintf("items.%d.location_end_id.unique", index), errUniqueLocation)
			}
		}
	}
}

func (oi *item) Save(so *model.StorageGroup) {
	o := orm.NewOrm()

	o.Raw("INSERT INTO storage_group_area (storage_group_id, warehouse_area_id, location_from_id, location_end_id) "+
		"VALUES (?,?,?,?);", so.ID, oi.Area.ID, oi.LocationFrom.ID, oi.LocationEnd.ID).Exec()

	var ls []*warehouse.Location
	o.Raw("SELECT * FROM warehouse_location where id >= ? and id <= ? and warehouse_area_id = ?", oi.LocationFrom.ID, oi.LocationEnd.ID, oi.Area.ID).QueryRows(&ls)

	if len(ls) > 0 {
		var vals []string
		for _, l := range ls {
			vals = append(vals, fmt.Sprintf("(%d,%d)", so.ID, l.ID))
		}

		o.Raw("INSERT INTO storage_group_location (storage_group_id, warehouse_location_id) VALUES " + strings.Join(vals, ",")).Exec()
	}
}
