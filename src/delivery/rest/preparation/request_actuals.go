// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"fmt"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/delivery/model"
)

type actual struct {
	ID   string `json:"id"`
	Note string `json:"note"`

	PreparationActual *model.PreparationActual `json:"-"`
}

func (rp *actual) Validate(index int, o *validation.Output) {
	var e error

	if rp.ID != "" {
		if rp.PreparationActual, e = validPreparationActual(rp.ID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.id.invalid", index), errInvalidPreparationActual)
		}
	}
}

func (rp *actual) Save() {
	rp.PreparationActual.Note = rp.Note
	rp.PreparationActual.Save("note")
}
