// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"strings"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/storage/model"
)

func createAreas(i *model.StorageGroup, attribs []area, delete bool) {
	o := orm.NewOrm()
	if delete {
		// menghapus semua attributes dari item ini
		o.Raw("DELETE FROM storage_group_area where storage_group_id = ?", i.ID).Exec()
	}

	if len(attribs) > 0 {
		var vals []string
		for _, a := range attribs {
			if id, e := common.Decrypt(a.AreaID); e == nil {
				vals = append(vals, fmt.Sprintf("(%d,%d)", i.ID, id))
			}
		}

		o.Raw("INSERT INTO storage_group_area (storage_group_id, warehouse_area_id) VALUES " + strings.Join(vals, ",")).Exec()
	}
}
