// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package opname

import (
	"time"

	"git.qasico.com/cuxs/cuxs/event"

	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
)

type commitRequest struct {
	ID int64 `json:"-"`

	StockOpname *model.StockOpname `json:"-"`
	Session     *auth.SessionData  `json:"-"`
}

func (ar *commitRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ar.StockOpname, e = validStockOpname(ar.ID); e != nil {
		o.Failure("id.invalid", errInvalidStockOpname)
	}

	return o
}

func (ar *commitRequest) Messages() map[string]string {
	return map[string]string{}
}

func (ar *commitRequest) Save() (e error) {
	ar.StockOpname.Status = "finish"
	ar.StockOpname.ApprovedBy = ar.Session.User.(*user.User)
	ar.StockOpname.ApprovedAt = time.Now()

	if e = ar.StockOpname.Save("status", "approved_by", "approved_at"); e == nil {
		commitItem(ar.StockOpname)

		go event.Call("stockopname::finished", ar.StockOpname)
	}

	return
}

func commitItem(so *model.StockOpname) {
	// ambil item stock opname
	o := orm.NewOrm()

	o.LoadRelated(so, "Items", 2)

	cont := make(map[int8]*model.StockStorage)
	for _, item := range so.Items {
		if item.IsVoid == 0 {
			var ss *model.StockStorage
			var exist bool
			if ss, exist = cont[item.ContrainerNum]; !exist {
				ss = &model.StockStorage{
					Container: item.Container,
					Location:  so.Location,
				}
				ss.Save()
				cont[item.ContrainerNum] = ss
			}

			// update stock unit ss
			item.Unit.Storage = ss
			item.Unit.Status = "stored"
			item.Unit.IsDefect = item.IsDefect
			item.Unit.Save("storage_id", "status", "is_defect")
		} else {
			if item.Unit.Storage != nil {
				if item.Unit.Storage.Location.ID == so.Location.ID {
					item.Unit.Storage = nil
					item.Unit.Status = "void"
					item.Unit.Save("storage_id", "status")
				}
			}
		}

		go event.Call("stockopname::commited", item)
	}
}
