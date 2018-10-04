// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/receiving/model"
)

// Show find a single data warehouse_area using field and value condition.
func Show(id int64) (*model.ReceivingUnit, error) {
	m := new(model.ReceivingUnit)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}
