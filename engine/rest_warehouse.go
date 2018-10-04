// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.qasico.com/gudang/api/src/warehouse/rest/area"
	"git.qasico.com/gudang/api/src/warehouse/rest/location"
)

func init() {
	handlers["warehouse/area"] = &area.Handler{}
	handlers["warehouse/location"] = &location.Handler{}
}
