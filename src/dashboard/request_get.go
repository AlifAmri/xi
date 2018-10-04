// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dashboard

import "git.qasico.com/cuxs/orm"

type TotalStorage struct {
	TotalLocation int64         `json:"total_location"`
	TotalCapacity int64         `json:"total_capacity"`
	TotalUsed     int64         `json:"total_used"`
	TotalStock    []*TotalStock `json:"total_stock"`
}

type TotalStock struct {
	Name  string  `json:"name"`
	Stock float64 `json:"stock"`
}

func show() (ts *TotalStorage, e error) {
	o := orm.NewOrm()

	if e = o.Raw("SELECT count(*) as total_location, sum(storage_capacity) as total_capacity, sum(storage_used) as total_used FROM warehouse_location").QueryRow(&ts); e == nil {
		var tts []*TotalStock
		o.Raw("SELECT it.name, sum(i.stock) as stock FROM item i " +
			"inner join item_type it on it.id = i.type_id group by i.type_id").QueryRows(&tts)

		ts.TotalStock = tts
	}

	return
}
