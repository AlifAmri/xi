// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"time"

	"git.qasico.com/cuxs/common/log"
	model3 "git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/receiving/model"

	inventory "git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/stock/model"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/stock/rest/unit"
)

type LogMaker struct {
	Doc       interface{}
	StockUnit *model.StockUnit
	Item      *inventory.Item
	Batch     *inventory.ItemBatch
	Quantity  float64
}

func CreateLog(lm *LogMaker) (sl *model.StockLog, e error) {
	var refCode string
	var refID uint64
	var refType string

	if d, ok := lm.Doc.(*model.StockOpname); ok {
		refCode = d.Code
		refID = uint64(d.ID)
		refType = "stock_adjustment"

		if d.Type == "opname" {
			refType = "stock_opname"
		}
	} else if d, ok := lm.Doc.(*model.StockMovement); ok {
		refCode = d.Code
		refID = uint64(d.ID)
		refType = "stock_movement"
	} else if d, ok := lm.Doc.(model2.Receiving); ok {
		refCode = d.DocumentCode
		refID = uint64(d.ID)
		refType = "receiving"
	} else if d, ok := lm.Doc.(*model3.Preparation); ok {
		refCode = d.DocumentCode
		refID = uint64(d.ID)
		refType = "preparation"
	}

	sl = &model.StockLog{
		StockUnit:  lm.StockUnit,
		Item:       lm.Item,
		Batch:      lm.Batch,
		Quantity:   lm.Quantity,
		RefType:    refType,
		RefCode:    refCode,
		RefID:      refID,
		RecordedAt: time.Now(),
	}

	if e = sl.Save(); e == nil {
		Recalculate(lm)
	}

	log.Error(e)

	return
}

func Recalculate(lm *LogMaker) {
	go func() {
		if lm.StockUnit != nil {
			calculateUnit(lm.StockUnit)
		}

		if lm.Batch != nil {
			calculateBatch(lm.Batch)
		}

		if lm.Item != nil {
			calculateItem(lm.Item)
		}
	}()
}

func RecalculateByStockUnit(rq *orm.RequestQuery) {
	go func() {
		sus, _, _ := unit.Get(rq)
		for _, su := range *sus {
			calculateUnit(&su)

			if su.Batch != nil {
				calculateBatch(su.Batch)
			}

			if su.Item != nil {
				calculateItem(su.Item)
			}
		}

	}()
}

func calculateUnit(su *model.StockUnit) {
	var stock float64
	orm.NewOrm().Raw("SELECT sum(quantity) FROM stock_log where stock_unit_id = ?;", su.ID).QueryRow(&stock)

	su.Stock = stock
	su.Save("stock")
}

func calculateBatch(ib *inventory.ItemBatch) {
	var stock float64
	orm.NewOrm().Raw("SELECT sum(quantity) FROM stock_log sl inner join stock_unit su on sl.stock_unit_id = su.id where sl.batch_id = ? AND su.status not in ('void', 'out');", ib.ID).QueryRow(&stock)

	ib.Stock = stock
	ib.Save("stock")
}

func calculateItem(i *inventory.Item) {
	var stock float64
	orm.NewOrm().Raw("SELECT sum(quantity) FROM stock_log sl inner join stock_unit su on sl.stock_unit_id = su.id where sl.item_id = ? AND su.status not in ('void', 'out');", i.ID).QueryRow(&stock)

	i.Stock = stock
	i.Save("stock")
}
