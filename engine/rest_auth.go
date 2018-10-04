// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	authRest "git.qasico.com/gudang/api/src/auth/rest"
	pRest "git.qasico.com/gudang/api/src/privilege/rest"
)

func init() {
	handlers["auth"] = &authRest.Handler{}
	handlers["privilege"] = &pRest.Handler{}
}
