// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"fmt"
	"git.qasico.com/cuxs/orm"
	model2 "git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/warehouse"
	"strings"
	"time"

	"git.qasico.com/cuxs/common/log"
	"git.qasico.com/cuxs/env"
	"github.com/tealeg/xlsx"
)

type BackdateStock struct {
	Item  int64
	Batch int64
	Stock float64
}

func stockItemBackDate(date time.Time, items []*model2.Item) {
	var ids []int64
	var tty []string
	for _, i := range items {
		ids = append(ids, i.ID)
		tty = append(tty, "?")
	}

	tandaTanya := strings.Join(tty, ",")

	var data []BackdateStock
	orm.NewOrm().Raw("SELECT item_id as item, sum(quantity) as stock FROM stock_log "+
		"where recorded_at <= ? and item_id in ("+tandaTanya+") group by item_id", date, ids).QueryRows(&data)

	for _, item := range items {
		item.Stock = 0
		for _, d := range data {
			if d.Item == item.ID {
				item.Stock = d.Stock
			}
		}
	}

	return
}

func stockBatchBackDate(date time.Time, batchs []*model2.ItemBatch) {
	var ids []int64
	var tty []string
	for _, i := range batchs {
		ids = append(ids, i.ID)
		tty = append(tty, "?")
	}

	tandaTanya := strings.Join(tty, ",")

	var data []BackdateStock
	orm.NewOrm().Raw("SELECT batch_id as batch, sum(quantity) as stock FROM stock_log "+
		"where recorded_at <= ? and batch_id in ("+tandaTanya+") group by batch_id", date, ids).QueryRows(&data)

	for _, item := range batchs {
		item.Stock = 0
		for _, d := range data {
			if d.Batch == item.ID {
				item.Stock = d.Stock
			}
		}
	}

	return
}

// GetUnit get all data item_type that matched with query request parameters.
// returning slices of ItemType, total data without limit and error.
func GetUnit(rq *orm.RequestQuery, wa *warehouse.Area) (m *[]model.StockUnit, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.StockUnit))

	q = q.Filter("storage__location__warehouse_area_id__id", wa.ID)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.StockUnit
	if _, err = q.RelatedSel().All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func GetUnitXls(wa *warehouse.Area, r []model.StockUnit) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")
	path := env.GetString("EXPORT_PATH", "")

	now := time.Now()

	name := strings.Replace(wa.Name, " ", "_", 1)
	filename := fmt.Sprintf("StockUnit-%s-%s.xlsx", name, now.Format("200601021504"))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", path, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = fmt.Sprintf("Data Stock Unit %s", wa.Name)
		row = sheet.AddRow()
		row.AddCell().Value = "Tanggal: " + now.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "KODE"
		row.AddCell().Value = "KODE ITEM"
		row.AddCell().Value = "KODE BATCH"
		row.AddCell().Value = "LOKASI"
		row.AddCell().Value = "QUANTITY"
		row.AddCell().Value = "STATUS"
		row.AddCell().Value = "NCP"

		for _, i := range r {
			row = sheet.AddRow()
			row.AddCell().Value = i.Code
			row.AddCell().Value = i.Item.Code
			row.AddCell().Value = i.Batch.Code
			row.AddCell().Value = i.Storage.Location.Name
			row.AddCell().SetFloat(i.Stock)
			row.AddCell().Value = i.Status
			if i.IsDefect == 1 {
				row.AddCell().Value = "NCP"
			}
		}

		err = file.Save(fileDir)
		log.Error(err)
	}

	return
}

func GetStockItemXls(date time.Time, r []*model2.Item) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")
	path := env.GetString("EXPORT_PATH", "")

	filename := fmt.Sprintf("Stock-%s.xlsx", date.Format("20060102"))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", path, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Data Stock"
		row = sheet.AddRow()
		row.AddCell().Value = "Tanggal: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "KODE"
		row.AddCell().Value = "QUANTITY"

		for _, i := range r {
			row = sheet.AddRow()
			row.AddCell().Value = i.Code
			row.AddCell().SetFloat(i.Stock)
		}

		err = file.Save(fileDir)
		log.Error(err)
	}

	return
}

func GetStockBatchXls(date time.Time, r []*model2.ItemBatch) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")
	path := env.GetString("EXPORT_PATH", "")

	filename := fmt.Sprintf("StockBatch-%s.xlsx", date.Format("20060102"))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", path, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Data Stock Batch"
		row = sheet.AddRow()
		row.AddCell().Value = "Tanggal: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "KODE ITEM"
		row.AddCell().Value = "KODE BATCH"
		row.AddCell().Value = "QUANTITY"

		for _, i := range r {
			row = sheet.AddRow()
			row.AddCell().Value = i.Item.Code
			row.AddCell().Value = i.Code
			row.AddCell().SetFloat(i.Stock)
		}

		err = file.Save(fileDir)
		log.Error(err)
	}

	return
}
