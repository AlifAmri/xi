// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package warehouse

import (
	"encoding/json"
	"fmt"
	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/common/log"
	"git.qasico.com/cuxs/env"
	"git.qasico.com/cuxs/orm"
	"github.com/tealeg/xlsx"
)

type ReportWarehouse struct {
	ID             int64   `orm:"column(id)" json:"-"`
	AreaName       string  `json:"area_name"`
	TotalLocation  int64   `json:"total_location"`
	TotalInactive  int64   `json:"total_inactive"`
	TotalCapacity  int64   `json:"total_capacity"`
	TotalUsed      int64   `json:"total_used"`
	TotalAvailable int64   `json:"total_available"`
	TotalQuantity  float64 `json:"total_quantity"`
}

func (m *ReportWarehouse) MarshalJSON() ([]byte, error) {
	type Alias ReportWarehouse

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

func Get() (m []ReportWarehouse, total int64, e error) {
	o := orm.NewOrm()
	total, e = o.Raw("SELECT wa.id, wa.name as area_name, count(wl.id) as total_location,sum(wl.storage_capacity) as total_capacity, sum(wl.storage_used) as total_used " +
		", (select count(wlx.id) from warehouse_location wlx where wlx.warehouse_area_id = wa.id and wlx.is_active = 0) as total_inactive " +
		", (sum(wl.storage_capacity) - sum(wl.storage_used) - (select IFNULL(sum(wlx.storage_capacity), 0) from warehouse_location wlx where wlx.warehouse_area_id = wa.id and wlx.is_active = 0)) as total_available " +
		", (SELECT sum(stock) FROM stock_unit su " +
		"inner join stock_storage st on st.id = su.storage_id " +
		"inner join warehouse_location stwl on stwl.id = st.location_id " +
		"where su.status in ('stored', 'moving', 'prepared') and stwl.warehouse_area_id = wa.id " +
		") as total_quantity " +
		"FROM warehouse_area wa left join warehouse_location wl on wl.warehouse_area_id = wa.id " +
		"group by wa.id;").QueryRows(&m)

	return
}

func Show(id int64) (m *ReportWarehouse, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT wa.id as id, wa.name as area_name, count(wl.id) as total_location,sum(wl.storage_capacity) as total_capacity, sum(wl.storage_used) as total_used "+
		", (select count(wlx.id) from warehouse_location wlx where wlx.warehouse_area_id = wa.id and wlx.is_active = 0) as total_inactive "+
		", (sum(wl.storage_capacity) - (select IFNULL(sum(wlx.storage_capacity), 0) from warehouse_location wlx where wlx.warehouse_area_id = wa.id and wlx.is_active = 0)) as total_available "+
		", (SELECT sum(stock) FROM stock_unit su "+
		"inner join stock_storage st on st.id = su.storage_id "+
		"inner join warehouse_location stwl on stwl.id = st.location_id "+
		"where su.status in ('stored', 'moving', 'prepared') and stwl.warehouse_area_id = wa.id "+
		") as total_quantity "+
		"FROM warehouse_area wa left join warehouse_location wl on wl.warehouse_area_id = wa.id "+
		"WHERE wa.id = ? "+
		"group by wa.id;", id).QueryRow(&m)

	return
}

func toXls(r []ReportWarehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")
	path := env.GetString("EXPORT_PATH", "")

	now := time.Now()
	filename := fmt.Sprintf("LaporanStorage-%s.xlsx", now.Format("200601021504"))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", path, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Laporan Storage"
		row = sheet.AddRow()
		row.AddCell().Value = "Tanggal: " + now.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Nama Area"
		row.AddCell().Value = "Jumlah Lokasi"
		row.AddCell().Value = "Total Kapasitas"
		row.AddCell().Value = "Jumlah Lokasi Tidak Aktif"
		row.AddCell().Value = "Kapasitas Tersedia"
		row.AddCell().Value = "Kapasitas Terpakai"
		row.AddCell().Value = "Total Quantity"

		for _, i := range r {
			row = sheet.AddRow()
			row.AddCell().Value = i.AreaName
			row.AddCell().SetInt64(i.TotalLocation)
			row.AddCell().SetInt64(i.TotalCapacity)
			row.AddCell().SetInt64(i.TotalInactive)
			row.AddCell().SetInt64(i.TotalAvailable)
			row.AddCell().SetInt64(i.TotalUsed)
			row.AddCell().SetFloat(i.TotalQuantity)
		}

		err = file.Save(fileDir)
		log.Error(err)
	}

	return
}
