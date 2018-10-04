// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	userRest "git.qasico.com/gudang/api/src/user/rest"
	usergroupRest "git.qasico.com/gudang/api/src/usergroup/rest"
)

func init() {
	handlers["usergroup"] = &usergroupRest.Handler{}
	handlers["user"] = &userRest.Handler{}
}
