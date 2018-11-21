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
	for _, d := range ru.Receiving.Actuals {
		if d.Batch != nil && d.Item.ID == ru.Unit.Item.ID && d.Batch.ID == ru.Unit.Batch.ID {
			if ru.IsNcp == 1 {
				d.QuantityDefect += ru.Quantity
			}

			d.QuantityReceived += ru.Quantity

			d.Save("quantity_defect", "quantity_received")
			updated = true
		}
	}

	if !updated {
		for _, d := range ru.Receiving.Actuals {
			if !updated && d.Batch == nil {
				if d.Item.ID == ru.Unit.Item.ID {
					if ru.IsNcp == 1 {
						d.QuantityDefect += ru.Quantity
					}

					d.QuantityReceived += ru.Quantity

					d.Save("quantity_defect", "quantity_received")
					updated = true
				}
			}
		}
	}

	if !updated {
		ra := &model.ReceivingActual{
			Receiving:        ru.Receiving,
			Item:             ru.Unit.Item,
			Batch:            ru.Unit.Batch,
			QuantityReceived: ru.Quantity,
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
	o.LoadRelated(r, "Documents", 0)

	for _, d := range r.Documents {
		ra := &model.ReceivingActual{
			Receiving:       r,
			Item:            d.Item,
			Batch:           d.Batch,
			QuantityPlanned: d.Quantity,
		}

		ra.Save()
	}

	o.LoadRelated(r, "Units", 1)

	for _, u := range r.Units {
		if u.Unit != nil {
			CalculateActualFromUnit(u)
		}
	}
}
