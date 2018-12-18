// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package services

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/receiving/model"
)

func CalculateActualFromUnit(ru *model.ReceivingUnit) {
	o := orm.NewOrm()

	o.LoadRelated(ru.Receiving, "Actuals", 1)

	var updated bool
	// loop actual untuk mencari item yang sama dengan receiving unit (item,batch dan unit yang sama)
	for _, d := range ru.Receiving.Actuals {
		if d.Batch != nil && d.Item.ID == ru.Unit.Item.ID && d.Batch.ID == ru.Unit.Batch.ID {
			if d.Unit != nil && d.Unit.Code == ru.UnitCode {
				if ru.IsNcp == 1 {
					d.QuantityDefect += ru.Quantity
				}

				d.QuantityReceived += ru.Quantity

				d.Save("quantity_defect", "quantity_received")
				updated = true
			}
		}
	}

	//if !updated {
	//	for _, d := range ru.Receiving.Actuals {
	//		if !updated && d.Batch == nil {
	//			if d.Item.ID == ru.Unit.Item.ID {
	//				if ru.IsNcp == 1 {
	//					d.QuantityDefect += ru.Quantity
	//				}
	//
	//				d.QuantityReceived += ru.Quantity
	//
	//				d.Save("quantity_defect", "quantity_received")
	//				updated = true
	//			}
	//		}
	//	}
	//}

	// jika tidak terdapat actual yang match dengan unit,
	if !updated {
		ra := &model.ReceivingActual{
			Receiving:        ru.Receiving,
			Item:             ru.Unit.Item,
			Batch:            ru.Unit.Batch,
			QuantityReceived: ru.Quantity,
			Unit:             ru.Unit,
		}
		if ru.IsNcp == 1 {
			ra.QuantityDefect = ru.Quantity
		}

		ra.Save()
	}
}

func CreateActual(r *model.Receiving) {
	o := orm.NewOrm()

	// hapus semua actual
	o.Raw("DELETE FROM receiving_actual where receiving_id = ?", r.ID).Exec()

	// recreate actual
	o.LoadRelated(r, "Documents", 1)

	for _, d := range r.Documents {
		ra := &model.ReceivingActual{
			Receiving:       r,
			Item:            d.Item,
			Batch:           d.Batch,
			QuantityPlanned: d.Quantity,
		}
		if d.Unit != nil {
			ra.Unit = d.Unit
		}
		ra.Save()
	}

	o.LoadRelated(r, "Units", 2)

	for _, u := range r.Units {
		if u.Unit != nil {
			CalculateActualFromUnit(u)
		}
	}
}
