// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"git.qasico.com/cuxs/common/log"
	"git.qasico.com/cuxs/cuxs"

	_ "git.qasico.com/gudang/api/custom/auth"
	"git.qasico.com/gudang/api/engine"
)

// init preparing application instances.
func init() {
	log.DebugMode = cuxs.IsDebug()
	log.Log = log.New()

	if e := cuxs.DbSetup(); e != nil {
		panic(e)
	}
}

// main creating new instances application
// and serving application server.
func main() {
	// starting server
	cuxs.StartServer(engine.Router())
}
