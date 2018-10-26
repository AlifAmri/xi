// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package item

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
)

// Show find a single data item using field and value condition.
func Show(id int64) (*model.Item, error) {
	m := new(model.Item)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter("id", id).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	var attrs []attribute
	if i, _ := o.Raw("SELECT attribute, value FROM item_attribute where item_id = ?", m.ID).QueryRows(&attrs); i > 0 {
		m.Attributes = make(map[string]string, i)

		for _, a := range attrs {
			m.Attributes[a.Attribute] = a.Value
		}
	}

	if m.Equation != "" {
		o.Raw("SELECT * FROM item where equation = ? and type_id = ? and id != ?", m.Equation, m.Type.ID, m.ID).QueryRows(&m.Equations)
	}

	return m, nil
}

// Get get all data item that matched with query request parameters.
// returning slices of Item, total data without limit and error.
func Get(rq *orm.RequestQuery) (m *[]*model.Item, total int64, err error) {
	// make new orm query
	q, o := rq.Query(new(model.Item))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Item
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, ig := range mx {
			var attrs []attribute

			if i, _ := o.Raw("SELECT attribute, value FROM item_attribute where item_id = ?", ig.ID).QueryRows(&attrs); i > 0 {
				ig.Attributes = make(map[string]string, i)
				for _, a := range attrs {
					ig.Attributes[a.Attribute] = a.Value
				}
			}
		}

		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
