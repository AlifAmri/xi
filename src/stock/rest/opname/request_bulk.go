// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package opname

import (
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"

	"errors"
	"fmt"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/inventory"
	"git.qasico.com/gudang/api/src/inventory/model"
	model2 "git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/warehouse"
)

type bulkRequest struct {
	Items   []*item           `json:"items" valid:"required"`
	Session *auth.SessionData `json:"-"`
	createR []*createRequest  `json:"-"`
}

type item struct {
	No           int64   `json:"no" valid:"required"`
	Location     string  `json:"location" valid:"required"`
	StockContent string  `json:"stock_content"`
	CodeSize     string  `json:"code_size" valid:"required"`
	Weekly       string  `json:"weekly" valid:"required"`
	PalletType   string  `json:"pallet_type" valid:"required"`
	NoPallet     string  `json:"no_pallet"`
	Quantity     float64 `json:"quantity" valid:"required"`
}

func (cr *bulkRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	loc := make(map[string][]*item)
	for _, m := range cr.Items {
		loc[m.Location] = append(loc[m.Location], m)
	}
	for locate, itm := range loc {
		var e error
		var req = &createRequest{Session: cr.Session}
		req.Note = ""
		if req.Location, e = validLocationBulk(locate); e != nil {
			o.Failure(fmt.Sprintf("location %s is", locate), "location is not valid")
		}

		if len(itm) > 0 {
			ctn := make(map[string]bool)
			for _, it := range itm {
				if it.StockContent != "" {
					if ctn[it.StockContent] {
						o.Failure(fmt.Sprintf("stock content %s dengan No %v tidak invalid", it.StockContent, it.No), "duplicate")
					} else {
						ctn[it.StockContent] = true
					}
				}
				itmop := validItemBulk(it, req.Location.ID, o)
				itmop.Note = ""

				req.Items = append(req.Items, itmop)
			}
		}
		cr.createR = append(cr.createR, req)
	}

	return o
}

func (cr *bulkRequest) Messages() map[string]string {
	return map[string]string{
		"items.required": errRequiredItems,
	}
}

func validLocationBulk(loc string) (l *warehouse.Location, e error) {
	l = new(warehouse.Location)
	l.Name = loc

	if e = l.Read("Name"); e == nil {
		if l.IsActive == 0 {
			e = errors.New(fmt.Sprintf("location %s is not active", loc))
		} else {
			// cek stockopname lain dengan location ini
			var total int64
			orm.NewOrm().Raw("SELECT count(id) from stock_opname where location_id = ? and status = ?", l.ID, "active").QueryRow(&total)
			if total > 0 {
				e = errors.New(fmt.Sprintf("location %s still have stockopname active", loc))
			}
		}
	}

	return
}

func validItemBulk(itm *item, locID int64, o *validation.Output) (ItemOp *opnameItem) {
	ItemOp = new(opnameItem)
	if itm.CodeSize != "" {
		ItemOp.Item = &model.Item{Code: itm.CodeSize, IsActive: 1, Type: &model.ItemType{ID: 1}}
		if e := ItemOp.Item.Read("Code", "IsActive", "Type"); e != nil {
			o.Failure(fmt.Sprintf("code size %s with No %v invalid", itm.CodeSize, itm.No), "tidak valid atau tidak aktif")
		} else {
			ItemOp.Item.Type.Read("ID")
			//cek stock content
			if itm.StockContent != "" {
				ItemOp.UnitCode = itm.StockContent
				ItemOp.StockUnit = &model2.StockUnit{Code: ItemOp.UnitCode}
				if e = ItemOp.StockUnit.Read("Code"); e == nil {
					if ItemOp.StockUnit.Storage != nil {
						ItemOp.StockUnit.Storage.Read("ID")
						if ItemOp.StockUnit.Storage.Location.ID != locID {
							o.Failure(fmt.Sprintf("lokasi stock content %s with No %v invalid", itm.StockContent, itm.No), "tidak valid atau lokasi stock content salah")
						}
					}
					ItemOp.StockUnit.Item.Read("ID")

					if ItemOp.StockUnit.Item.ID != ItemOp.Item.ID {
						o.Failure(fmt.Sprintf("stock content %s with No %v invalid", itm.StockContent, itm.No), "tidak valid atau sudah ada di item lain")
					}

				} else {
					ItemOp.StockUnit = nil
				}
			}

			if itm.Weekly != "" {
				if ItemOp.BatchCode, e = validBatchCode(itm.Weekly); e != nil {
					o.Failure(fmt.Sprintf("weekly %s with No %v invalid", itm.Weekly, itm.No), "format weekly salah")
				} else {
					ItemOp.ItemBatch = inventory.GetBatch(ItemOp.Item.ID, ItemOp.BatchCode)
				}
			}
			ItemOp.Quantity = itm.Quantity
			if itm.Quantity == 0 {
				o.Failure(fmt.Sprintf("quantity No %v salah", itm.No), "harus diisi")
			}
			if ItemOp.Container, e = validContainerBulk(itm.PalletType); e != nil {
				o.Failure(fmt.Sprintf("tipe pallet %s with No %v invalid", itm.PalletType, itm.No), "tipe pallet salah")
			}
			ItemOp.ContainerNum = int8(common.ToInt(itm.NoPallet))

		}
	}
	return
}

func (cr *bulkRequest) Save() (e error) {
	for _, data := range cr.createR {
		_, e = data.Save()
	}
	return
}
