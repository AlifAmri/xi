// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.qasico.com/gudang/api/src/partnership/rest/partner"
	"git.qasico.com/gudang/api/src/partnership/rest/type"
)

func init() {
	handlers["partnership/partner"] = &partner.Handler{}
	handlers["partnership/type"] = &ptype.Handler{}
}
