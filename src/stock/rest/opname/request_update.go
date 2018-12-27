// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package opname

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
)

type updateRequest struct {
	ID    int64         `json:"-" valid:"required"`
	Note  string        `json:"note"`
	Items []*opnameItem `json:"items" valid:"required"`

	Session     *auth.SessionData  `json:"-"`
	StockOpname *model.StockOpname `json:"-"`
}

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.StockOpname, e = validStockOpname(ur.ID); e != nil {
		o.Failure("id.invalid", errInvalidStockOpname)
	}

	if len(ur.Items) > 0 {
		contcollector := make(map[int8]bool)
		for i, item := range ur.Items {
			item.Validate(i, o)
			if item.IsVoid == int8(0) {
				contcollector[item.ContainerNum] = true
			}
		}
		var countMove int
		if ur.StockOpname != nil {
			ur.StockOpname.Location.Read("ID")
			countMove = countMovement(ur.StockOpname.Location.ID)
		}
		// cek max pallet dan container
		if (len(contcollector) + countMove) > ur.StockOpname.Location.StorageCapacity {
			o.Failure("location_id.invalid", "pallet di lokasi ini sudah maksimum")
		}
	}

	return o
}

func (ur *updateRequest) Messages() map[string]string {
	return map[string]string{
		"items.required": errRequiredItems,
	}
}

func (ur *updateRequest) Save() (u *model.StockOpname, e error) {
	ur.StockOpname.Note = ur.Note

	if e = ur.StockOpname.Save("note"); e == nil {
		for _, item := range ur.Items {
			item.Save(ur.StockOpname, ur.Session.User.(*user.User))
		}
	}

	return ur.StockOpname, e
}
