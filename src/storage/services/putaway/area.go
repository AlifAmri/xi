// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package storage

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/warehouse"
)

func findAreaNCP() (wa []*warehouse.Area) {
	orm.NewOrm().Raw("SELECT wa.* FROM storage_group_area sga " +
		"inner join warehouse_area wa on wa.id = sga.warehouse_area_id " +
		"inner join storage_group sg on sg.id = sga.storage_group_id " +
		"where sg.type = 'ncp' and sg.is_active = 1").QueryRows(&wa)

	return
}

func findAreaItem(code string) (wa []*warehouse.Area) {
	orm.NewOrm().Raw("SELECT wa.* FROM item i "+
		"inner join warehouse_area wa on wa.id = i.preferred_area "+
		"where i.code = ?;", code).QueryRows(&wa)

	return
}

func findAreaConfig(icode string, bcode string) (wa []*warehouse.Area) {
	if wa = configCode(icode); len(wa) > 0 {
		return
	} else if wa = configBatch(bcode); len(wa) > 0 {
		return
	} else if wa = configCategory(icode); len(wa) > 0 {
		return
	} else if wa = configGroup(icode); len(wa) > 0 {
		return
	}

	return configDefault()
}

func configCode(code string) (wa []*warehouse.Area) {
	c := code[0:3]
	orm.NewOrm().Raw("SELECT wa.* FROM storage_group_area sga "+
		"inner join warehouse_area wa on wa.id = sga.warehouse_area_id "+
		"inner join storage_group sg on sg.id = sga.storage_group_id "+
		"where sg.type = 'item_code' and sg.type_value = ? and sg.is_active = 1", c).QueryRows(&wa)

	return
}

func configCategory(code string) (wa []*warehouse.Area) {
	orm.NewOrm().Raw("SELECT wa.* FROM storage_group_area sga "+
		"inner join warehouse_area wa on wa.id = sga.warehouse_area_id "+
		"inner join storage_group sg on sg.id = sga.storage_group_id "+
		"where sg.type = 'item_code' and sg.is_active = 1 and sg.type_value = ("+
		"SELECT ic.name from item i "+
		"inner join item_category ic on ic.id = i.category_id "+
		"where i.code = ?)", code).QueryRows(&wa)

	return
}

func configGroup(code string) (wa []*warehouse.Area) {
	orm.NewOrm().Raw("SELECT wa.* FROM storage_group_area sga "+
		"inner join warehouse_area wa on wa.id = sga.warehouse_area_id "+
		"inner join storage_group sg on sg.id = sga.storage_group_id "+
		"where sg.type = 'item_code' and sg.is_active = 1  and sg.type_value = ("+
		"SELECT ig.name from item i "+
		"inner join item_group ig on ig.id = i.group_id "+
		"where i.code = ?)", code).QueryRows(&wa)

	return
}

func configBatch(code string) (wa []*warehouse.Area) {
	c := code[len(code)-2:]
	orm.NewOrm().Raw("SELECT wa.* FROM storage_group_area sga "+
		"inner join warehouse_area wa on wa.id = sga.warehouse_area_id "+
		"inner join storage_group sg on sg.id = sga.storage_group_id "+
		"where sg.type = 'item_code' and sg.type_value = ? and sg.is_active = 1", c).QueryRows(&wa)

	return
}

func configDefault() (wa []*warehouse.Area) {
	orm.NewOrm().Raw("SELECT wa.* FROM storage_group_area sga " +
		"inner join warehouse_area wa on wa.id = sga.warehouse_area_id " +
		"inner join storage_group sg on sg.id = sga.storage_group_id " +
		"where sg.type = 'default' and sg.is_active = 1").QueryRows(&wa)

	return
}
