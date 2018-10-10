// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.qasico.com/gudang/api/src/delivery/rest/order"
	"git.qasico.com/gudang/api/src/delivery/rest/plan"
	"git.qasico.com/gudang/api/src/delivery/rest/preparation"
	"git.qasico.com/gudang/api/src/delivery/rest/unit"
)

func init() {
	handlers["delivery/plan"] = &plan.Handler{}
	handlers["delivery/preparation"] = &preparation.Handler{}
	handlers["delivery/order"] = &order.Handler{}
	handlers["delivery/unit"] = &unit.Handler{}
}
