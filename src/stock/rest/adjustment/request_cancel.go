// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package adjustment

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/stock/model"
)

type cancelRequest struct {
	ID int64 `json:"-"`

	StockOpname *model.StockOpname `json:"-"`
	Session     *auth.SessionData  `json:"-"`
}

func (ar *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ar.StockOpname, e = validStockOpname(ar.ID); e != nil {
		o.Failure("id.invalid", errInvalidStockOpname)
	}

	return o
}

func (ar *cancelRequest) Messages() map[string]string {
	return map[string]string{}
}

func (ar *cancelRequest) Save() (e error) {
	ar.StockOpname.Status = "cancelled"

	e = ar.StockOpname.Save("status")

	return
}
