// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package services

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/delivery/model"
)

func CalculateActualFromUnit(ru *model.PreparationUnit) {
	ru.Unit.Read()
	o := orm.NewOrm()

	o.LoadRelated(ru.Preparation, "Actuals", 1)

	var updated bool
	for _, d := range ru.Preparation.Actuals {
		if d.Batch != nil && d.Item.ID == ru.Unit.Item.ID && d.Batch.ID == ru.Unit.Batch.ID {
			d.QuantityPrepared += ru.Quantity

			d.Save("quantity_prepared")
			updated = true
		}
	}

	if !updated {
		for _, d := range ru.Preparation.Actuals {
			if !updated && d.Item.ID == ru.Unit.Item.ID {
				d.QuantityPrepared += ru.Quantity

				d.Save("quantity_prepared")
				updated = true
			}
		}
	}

	if !updated {
		ra := &model.PreparationActual{
			Preparation:      ru.Preparation,
			Item:             ru.Unit.Item,
			Batch:            ru.Unit.Batch,
			QuantityPrepared: ru.Quantity,
		}

		ra.Save()
	}
}

func CreateActual(r *model.Preparation) {
	o := orm.NewOrm()

	// hapus semua actual
	o.Raw("DELETE FROM preparation_actual where preparation_id = ?", r.ID).Exec()

	// recreate actual
	o.LoadRelated(r, "Documents", 0)

	for _, d := range r.Documents {
		ra := &model.PreparationActual{
			Preparation:     r,
			Item:            d.Item,
			Batch:           d.Batch,
			QuantityPlanned: d.Quantity,
			Year:            d.Year,
		}

		ra.Save()
	}

	o.LoadRelated(r, "Units", 1)

	for _, u := range r.Units {
		CalculateActualFromUnit(u)
	}
}

func CalculateQuantity(r *model.Preparation) {
	var actual, planned float64
	o := orm.NewOrm()
	o.Raw("SELECT sum(quantity_prepared) as quantity FROM preparation_actual where preparation_id = ?;", r.ID).QueryRow(&actual)
	o.Raw("SELECT sum(quantity_planned) as quantity FROM preparation_actual where preparation_id = ?;", r.ID).QueryRow(&planned)

	r.TotalQuantityActual = actual
	r.TotalQuantityPlan = planned

	r.Save("total_quantity_actual", "total_quantity_plan")
}
