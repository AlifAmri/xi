// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"fmt"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/warehouse"
)

type SuggestionLocation struct {
	LocationKey  int64   `json:"location_id"`
	LocationName string  `json:"location_name"`
	Quantity     float64 `json:"quantity"`
}

// Show find a single data warehouse_area using field and value condition.
func Show(id int64) (*model.Preparation, error) {
	m := new(model.Preparation)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "Documents", 2)
	o.LoadRelated(m, "Units", 2)
	o.LoadRelated(m, "Actuals", 2)

	getSuggestion(m)

	return m, nil
}

// Get get all data warehouse_area that matched with query request parameters.
// returning slices of Area, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]model.Preparation, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.Preparation))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.Preparation
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func getSuggestion(p *model.Preparation) {

	// loop through documents
	for _, i := range p.Documents {
		var withBatch string
		var qp float64

		if i.Batch != nil {
			withBatch = fmt.Sprintf(" and su.batch_id = %d", i.Batch.ID)
		}

		orm.NewOrm().Raw("SELECT sum(pu.quantity) as quantity FROM preparation_unit pu "+
			"inner join stock_unit su on su.id = pu.unit_id "+
			"where su.item_id = ? and preparation_id = ?"+withBatch, i.Item.ID, p.ID).QueryRow(&qp)

		if r, e := getLocation(i.Item, i.Batch, i.Quantity-qp); e == nil {
			p.Pickings = append(p.Pickings, r...)
		}
	}
}

func getLocation(i *model2.Item, ib *model2.ItemBatch, q float64) (r []*model.PreparationSuggested, e error) {
	var withBatch string
	var totalQuantityChoosed float64

	if ib != nil {
		withBatch = fmt.Sprintf("and su.batch_id = %d ", ib.ID)
	}

	var s []*SuggestionLocation
	_, e = orm.NewOrm().Raw("SELECT max(wl.id) as location_key, max(wl.name) as location_name, sum(su.stock) as quantity FROM stock_unit su "+
		"inner join stock_storage ss on ss.id = su.storage_id "+
		"inner join warehouse_location wl on wl.id = ss.location_id "+
		"inner join warehouse_area wa on wa.id = wl.warehouse_area_id "+
		"where wa.type = 'storage' and su.status = 'stored' and su.is_defect = 0 and su.item_id = ? "+withBatch+
		"group by ss.location_id order by wl.warehouse_area_id DESC, wl.id ASC, quantity DESC;", i.ID).QueryRows(&s)

	if e == nil && len(s) > 0 {
		for _, si := range s {
			if totalQuantityChoosed < q {
				ps := &model.PreparationSuggested{
					Location:  &warehouse.Location{ID: si.LocationKey, Name: si.LocationName},
					Item:      i,
					ItemBatch: ib,
					Quantity:  si.Quantity,
				}

				if totalQuantityChoosed+si.Quantity > q {
					ps.Quantity = q - totalQuantityChoosed
				}

				totalQuantityChoosed += si.Quantity
				r = append(r, ps)
			}

			if totalQuantityChoosed == q {
				break
			}
		}
	}

	return
}
