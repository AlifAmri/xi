// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import "git.qasico.com/cuxs/orm"

func get() (v map[string]string, e error) {
	var attrs []attribute

	v = make(map[string]string)
	if _, e = orm.NewOrm().Raw("SELECT attribute, value FROM app_config").QueryRows(&attrs); e == nil {
		for _, a := range attrs {
			v[a.Attribute] = a.Value
		}
	}

	return
}
