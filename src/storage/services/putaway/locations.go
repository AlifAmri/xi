// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package putaway

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/warehouse"
)

func findLocationNCP() (wl []*warehouse.Location) {
	orm.NewOrm().Raw("SELECT wl.* FROM storage_group_location sgl " +
		"inner join warehouse_location wl on wl.id = sgl.warehouse_location_id " +
		"inner join storage_group sg on sg.id = sgl.storage_group_id " +
		"where sg.type = 'ncp' and sg.is_active = 1 " +
		"and wl.is_active = 1 and wl.storage_capacity > wl.storage_used " +
		"order by wl.storage_used DESC, wl.id ASC").QueryRows(&wl)

	return
}

func findLocationConfig(icode string, bcode string) (wl []*warehouse.Location) {
	if wls := configBatch(bcode); len(wls) > 0 {
		wl = append(wl, wls...)
	}

	if wls := configCode(icode); len(wls) > 0 {
		wl = append(wl, wls...)
	}

	if wls := configCategory(icode); len(wls) > 0 {
		wl = append(wl, wls...)
	}

	if wls := configGroup(icode); len(wls) > 0 {
		wl = append(wl, wls...)
	}

	if wls := configDefault(); len(wls) > 0 {
		wl = append(wl, wls...)
	}

	if wls := allLocation(); len(wls) > 0 {
		wl = append(wl, wls...)
	}

	return
}

func configCode(code string) (wl []*warehouse.Location) {
	c := code[0:3]
	orm.NewOrm().Raw("SELECT wl.* FROM storage_group_location sgl "+
		"inner join warehouse_location wl on wl.id = sgl.warehouse_location_id "+
		"inner join storage_group sg on sg.id = sgl.storage_group_id "+
		"where sg.type = 'item_code' and sg.type_value = ? and sg.is_active = 1 "+
		"and wl.is_active = 1 and wl.storage_capacity > wl.storage_used "+
		"order by wl.storage_used DESC, wl.id ASC", c).QueryRows(&wl)

	return
}

func configBatch(code string) (wl []*warehouse.Location) {
	c := code[len(code)-2:]
	orm.NewOrm().Raw("SELECT wl.* FROM storage_group_location sgl "+
		"inner join warehouse_location wl on wl.id = sgl.warehouse_location_id "+
		"inner join storage_group sg on sg.id = sgl.storage_group_id "+
		"where sg.type = 'item_batch' and sg.type_value = ? and sg.is_active = 1 "+
		"and wl.is_active = 1 and wl.storage_capacity > wl.storage_used "+
		"order by wl.storage_used DESC, wl.id ASC", c).QueryRows(&wl)

	return
}

func configCategory(code string) (wl []*warehouse.Location) {
	orm.NewOrm().Raw("SELECT wl.* FROM storage_group_location sgl "+
		"inner join warehouse_location wl on wl.id = sgl.warehouse_location_id "+
		"inner join storage_group sg on sg.id = sgl.storage_group_id "+
		"where sg.type = 'item_category' and sg.is_active = 1 and sg.type_value = ("+
		"SELECT ic.name from item i "+
		"inner join item_category ic on ic.id = i.category_id "+
		"where i.code = ?) "+
		"and wl.is_active = 1 and wl.storage_capacity > wl.storage_used "+
		"order by wl.storage_used DESC, wl.id ASC", code).QueryRows(&wl)

	return
}

func configGroup(code string) (wl []*warehouse.Location) {
	orm.NewOrm().Raw("SELECT wl.* FROM storage_group_location sgl "+
		"inner join warehouse_location wl on wl.id = sgl.warehouse_location_id "+
		"inner join storage_group sg on sg.id = sgl.storage_group_id "+
		"where sg.type = 'item_group' and sg.is_active = 1 and sg.type_value = ("+
		"SELECT ig.name from item i "+
		"inner join item_group ig on ig.id = i.group_id "+
		"where i.code = ?) "+
		"and wl.is_active = 1 and wl.storage_capacity > wl.storage_used "+
		"order by wl.storage_used DESC, wl.id ASC", code).QueryRows(&wl)

	return
}

func configDefault() (wl []*warehouse.Location) {
	orm.NewOrm().Raw("SELECT wl.* FROM storage_group_location sgl " +
		"inner join warehouse_location wl on wl.id = sgl.warehouse_location_id " +
		"inner join storage_group sg on sg.id = sgl.storage_group_id " +
		"where sg.type = 'default' and sg.is_active = 1 " +
		"and wl.is_active = 1 and wl.storage_capacity > wl.storage_used " +
		"order by wl.storage_used DESC, wl.id ASC").QueryRows(&wl)

	return
}

func allLocation() (wl []*warehouse.Location) {
	orm.NewOrm().Raw("SELECT * FROM warehouse_location wl " +
		"where wl.is_active = 1 and wl.storage_capacity > wl.storage_used " +
		"order by wl.storage_used DESC, wl.id ASC").QueryRows(&wl)

	return
}
