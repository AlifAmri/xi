// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.qasico.com/gudang/api/src/receiving/rest/plan"
	"git.qasico.com/gudang/api/src/receiving/rest/receiving"
	"git.qasico.com/gudang/api/src/receiving/rest/unit"
)

func init() {
	handlers["receiving/plan"] = &plan.Handler{}
	handlers["receiving/unit"] = &unit.Handler{}
	handlers["receiving"] = &receiving.Handler{}
}
