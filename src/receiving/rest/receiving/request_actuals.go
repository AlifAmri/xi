// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"fmt"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/receiving/model"
)

type actual struct {
	ID   string `json:"id"`
	Note string `json:"note"`

	ReceivingActual *model.ReceivingActual `json:"-"`
}

func (rp *actual) Validate(index int, o *validation.Output) {
	var e error

	if rp.ID != "" {
		if rp.ReceivingActual, e = validReceivingActual(rp.ID); e != nil {
			o.Failure(fmt.Sprintf("items.%d.id.invalid", index), errInvalidReceivingActual)
		}
	}
}

func (rp *actual) Save() {
	rp.ReceivingActual.Note = rp.Note
	rp.ReceivingActual.Save("note")
}
