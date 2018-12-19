// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"fmt"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/delivery/model"
)

// Show find a single data warehouse_area using field and value condition.
func Show(id int64) (*model.DeliveryOrder, error) {
	m := new(model.DeliveryOrder)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "Items", 1)
	return m, nil
}

// ShowPrint find a single data delivery_order using field and value condition with detailed data of preparation and collection of item type
func ShowPrint(id int64) (*model.DeliveryOrder, error) {
	m := new(model.DeliveryOrder)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	getLogo(m)
	counterPrint(m)
	o.LoadRelated(m, "Items", 1)
	tireGroup(m)
	return m, nil
}

func getLogo(m *model.DeliveryOrder){
	orm.NewOrm().Raw("select value from app_config where attribute = ?","logo").QueryRow(&m.Logo)
}

// counterPrint create counter for numeric print delivery order, each print will do increment for counter
func counterPrint(m *model.DeliveryOrder) {
	var id, value int64
	o := orm.NewOrm()
	e := o.Raw("SELECT c.id , c.value FROM counter c WHERE c.table_name = ?", "delivery_order_print").QueryRow(&id, &value)
	if e != nil || id == int64(0) {
		o.Raw("INSERT INTO counter (table_name, value, note) VALUES (?,?,'counter for DO print');", "delivery_order_print", int64(1)).Exec()
		m.Counter = fmt.Sprintf("%05d", 1)
	} else {
		value += int64(1)
		m.Counter = fmt.Sprintf("%05d", value)
		o.Raw("UPDATE counter SET value = ? WHERE id = ?", value, id).Exec()
	}
}

// tireGroup group quantity of tire based on type (4W and 2W)
func tireGroup(m *model.DeliveryOrder) {
	o := orm.NewOrm()
	for _, prepar := range m.Items {
		var quantity1, quantity2 float64
		o.Raw("SELECT SUM(pu.quantity) AS quantity1 FROM preparation_unit pu "+
			"INNER JOIN stock_unit su ON su.id = pu.unit_id "+
			"INNER JOIN item i ON i.id = su.item_id "+
			"INNER JOIN item_category ic ON ic.id = i.category_id "+
			"WHERE pu.preparation_id = ? AND pu.is_active = 1 AND ic.type_id = ? AND (ic.name = ? OR ic.name = ? );", prepar.ID, 1, "MC BIG", "MC SMALL").QueryRow(&quantity1)

		o.Raw("SELECT SUM(pu.quantity) AS quantity2 FROM preparation_unit pu "+
			"INNER JOIN stock_unit su ON su.id = pu.unit_id "+
			"INNER JOIN item i ON i.id = su.item_id "+
			"INNER JOIN item_category ic ON ic.id = i.category_id "+
			"WHERE pu.preparation_id = ? AND pu.is_active = 1 AND ic.type_id = ? AND ic.name != ? AND ic.name != ? ;", prepar.ID, 1, "MC BIG", "MC SMALL").QueryRow(&quantity2)
		m.Tire2W += quantity1
		m.Tire4W += quantity2
	}
}

// Get get all data warehouse_area that matched with query request parameters.
// returning slices of Area, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]model.DeliveryOrder, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.DeliveryOrder))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.DeliveryOrder
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
