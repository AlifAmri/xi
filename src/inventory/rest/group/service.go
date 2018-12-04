// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package group

import (
	"fmt"
	"strings"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
)

func createAttribute(ig *model.ItemGroup, attribs []attribute, delete bool) {
	o := orm.NewOrm()
	if delete {
		// menghapus semua attributes dari item group ini
		o.Raw("DELETE FROM item_group_attribute where item_group_id = ?", ig.ID).Exec()
	}

	if len(attribs) > 0 {
		var vals []string
		for _, a := range attribs {
			vals = append(vals, fmt.Sprintf("(%d,'%s','%s')", ig.ID, a.Attribute, a.Value))
		}

		o.Raw("INSERT INTO item_group_attribute (item_group_id, attribute, value) VALUES " + strings.Join(vals, ",")).Exec()
	}
}
