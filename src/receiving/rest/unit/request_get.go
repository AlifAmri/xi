// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/storage/services/putaway"
)

// Show find a single data warehouse_area using field and value condition.
func Show(id int64) (*model.ReceivingUnit, error) {
	m := new(model.ReceivingUnit)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	if m.LocationSuggested == nil {
		m.LocationSuggested = putaway.SuggestedPutaway(m.ItemCode, m.BatchCode, m.IsNcp)
		m.Save("location_suggested")
	}

	return m, nil
}
