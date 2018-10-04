// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/warehouse"
)

func init() {
	orm.RegisterModel(new(StockStorage))
}

// StockStorage model for stock_unit table.
type StockStorage struct {
	ID        int64               `orm:"column(id);auto" json:"-"`
	Location  *warehouse.Location `orm:"column(location_id);null;rel(fk)" json:"location,omitempty"`
	Container *model.Item         `orm:"column(container_id);null;rel(fk)" json:"container"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *StockStorage) MarshalJSON() ([]byte, error) {
	type Alias StockStorage

	alias := &struct {
		ID          string `json:"id"`
		ContainerID string `json:"container_id"`
		LocationID  string `json:"location_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	if m.Container != nil && m.Container.ID != int64(0) {
		alias.ContainerID = common.Encrypt(m.Container.ID)
	} else {
		alias.Container = nil
	}

	if m.Location != nil && m.Location.ID != int64(0) {
		alias.LocationID = common.Encrypt(m.Location.ID)
	} else {
		alias.Location = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating StockUnit struct into stock_unit table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to stock_unit.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *StockStorage) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting stock_unit data
// this also will truncated all data from all table
// that have relation with this stock_unit.
func (m *StockStorage) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}

// Read execute select based on data struct that already
// assigned.
func (m *StockStorage) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
