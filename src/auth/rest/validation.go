// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import "git.qasico.com/gudang/api/src/auth"

func validUsername(username string) (auth.UserModelInterface, error) {
	return auth.Service.GetByUsername(username)
}
