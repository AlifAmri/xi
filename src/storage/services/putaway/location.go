// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package storage

import "git.qasico.com/gudang/api/src/warehouse"

func getLocation(icode string, bcode string, wa []*warehouse.Area) (wl *warehouse.Location) {
	// cari lokasi yang sama item & batch di area ini
	// order berdasarkan storage_capacity - storage_used
	// kalau tidak ada cari item yang sama
	// order berdasarkan storage_capacity - storage_used
	// kalau nggak ada cari lokasi yang kosong
	// kalau nggak ada bebas
}
