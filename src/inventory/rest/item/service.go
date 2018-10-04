// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package item

import (
	"fmt"
	"strings"

	"git.qasico.com/gudang/api/src/inventory/model"

	"git.qasico.com/cuxs/orm"
)

func createAttribute(i *model.Item, attribs []attribute, delete bool) {
	o := orm.NewOrm()
	if delete {
		// menghapus semua attributes dari item ini
		o.Raw("DELETE FROM item_attribute where item_id = ?", i.ID).Exec()
	}

	if len(attribs) > 0 {
		var vals []string
		for _, a := range attribs {
			vals = append(vals, fmt.Sprintf("(%d,'%s','%s')", i.ID, a.Attribute, a.Value))
		}

		o.Raw("INSERT INTO item_attribute (item_id, attribute, value) VALUES " + strings.Join(vals, ",")).Exec()
	}
}
