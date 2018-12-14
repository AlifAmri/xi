// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/receiving/model"
)

type deleteRequest struct {
	ID int64 `json:"-" valid:"required"`

	ReceivingUnit *model.ReceivingUnit `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

func (cr *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.ReceivingUnit, e = validReceivingUnit(cr.ID, "draft"); e != nil {
		o.Failure("id.invalid", errInvalidReceivingUnit)
	}

	return o
}

func (cr *deleteRequest) Messages() map[string]string {
	return map[string]string{}
}

func (cr *deleteRequest) Save() (e error) {
	e = cr.ReceivingUnit.Delete()
	// hapus stock unit jika tipe draft dan tidak ada di receiving document
	var tot int64
	or := orm.NewOrm()
	or.Raw("SELECT count(*) FROM receiving_document rd "+
		"INNER JOIN receiving r ON r.id = rd.receiving_id "+
		"INNER JOIN stock_unit su ON su.id = rd.unit_id "+
		"WHERE su.code = ? AND r.id = ?", cr.ReceivingUnit.UnitCode, cr.ReceivingUnit.Receiving.ID).QueryRow(&tot)
	if tot == int64(0) {
		or.Raw("DELETE FROM stock_unit  WHERE code = ? AND status = ?", cr.ReceivingUnit.UnitCode, "draft").Exec()
	}
	return
}
