// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package adjustment

import (
	"time"

	"git.qasico.com/cuxs/cuxs/event"
	"git.qasico.com/cuxs/orm"

	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"

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
	}

	return
}

func commitItem(so *model.StockOpname) {
	// ambil item stock opname
	o := orm.NewOrm()

	o.LoadRelated(so, "Items", 2)

	for _, item := range so.Items {
		go event.Call("stockopname::commited", item)
	}
}
