// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"fmt"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/delivery/model"
)

type item struct {
	PreparationID string `json:"preparation_id" valid:"required"`

	Preparation *model.Preparation `json:"-"`
}

func (rp *item) Validate(index int, o *validation.Output, do *model.DeliveryOrder) {
	var e error

	if rp.PreparationID != "" {
		if rp.Preparation, e = validPreparation(rp.PreparationID, do); e != nil {
			o.Failure(fmt.Sprintf("items.%d.preparation_id.invalid", index), errInvalidPreparation)
		}
	}
}

func (rp *item) Save(r *model.DeliveryOrder) {
	rp.Preparation.DeliveryOrder = r
	rp.Preparation.Save("delivery_order_id")
}
