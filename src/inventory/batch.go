// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package inventory

import "git.qasico.com/gudang/api/src/inventory/model"

// GetBatch mendapatkan data item batch
// kalau tidak ada, fungsi ini akan membuat item batch baru
func GetBatch(itemID int64, code string) *model.ItemBatch {
	ib := &model.ItemBatch{Item: &model.Item{ID: itemID}, Code: code}

	if e := ib.Read("item_id", "code"); e != nil {
		ib.Save()
	}

	return ib
}

// GetItem mendapatkan data item
// kalau tidak ada, fungsi ini akan membuat item baru
func GetItem(code string) *model.Item {
	i := &model.Item{Type: &model.ItemType{ID: 1}, Code: code, IsActive: 1}

	if e := i.Read("type_id", "code"); e != nil {
		i.Save()
	}

	return i
}
