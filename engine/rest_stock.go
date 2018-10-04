// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.qasico.com/gudang/api/src/stock/rest/adjustment"
	"git.qasico.com/gudang/api/src/stock/rest/log"
	"git.qasico.com/gudang/api/src/stock/rest/movement"
	"git.qasico.com/gudang/api/src/stock/rest/opname"
	"git.qasico.com/gudang/api/src/stock/rest/unit"
)

func init() {
	handlers["stock/unit"] = &unit.Handler{}
	handlers["stock/opname"] = &opname.Handler{}
	handlers["stock/adjust"] = &adjustment.Handler{}
	handlers["stock/movement"] = &movement.Handler{}
	handlers["stock/log"] = &log.Handler{}
}
