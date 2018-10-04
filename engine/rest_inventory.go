// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.qasico.com/gudang/api/src/inventory/rest/batch"
	"git.qasico.com/gudang/api/src/inventory/rest/category"
	"git.qasico.com/gudang/api/src/inventory/rest/group"
	"git.qasico.com/gudang/api/src/inventory/rest/item"
	"git.qasico.com/gudang/api/src/inventory/rest/item_type"
	"git.qasico.com/gudang/api/src/inventory/rest/uom"
)

func init() {
	handlers["inventory/category"] = &category.Handler{}
	handlers["inventory/group"] = &group.Handler{}
	handlers["inventory/item"] = &item.Handler{}
	handlers["inventory/type"] = &itype.Handler{}
	handlers["inventory/uom"] = &uom.Handler{}
	handlers["inventory/batch"] = &batch.Handler{}
}
