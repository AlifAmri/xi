// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package movement

import (
	"time"

	"git.qasico.com/cuxs/cuxs/event"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/gudang/api/src/auth"
	model2 "git.qasico.com/gudang/api/src/inventory/model"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/warehouse"
)

type finishRequest struct {
	ID          int64  `json:"-" valid:"required"`
	StorageID   string `json:"storage_id"`
	ContainerID string `json:"container_id" valid:"required"`

	StockMovement *model.StockMovement `json:"-"`
	Session       *auth.SessionData    `json:"-"`
	Storage       *model.StockStorage  `json:"-"`
	Container     *model2.Item         `json:"-"`
}

func (cr *finishRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.StockMovement, e = validStockMovement(cr.ID, "start"); e != nil {
		o.Failure("id.invalid", errInvalidStockMovement)
	}

	if cr.StorageID != "" {
		if cr.Storage, e = validStorage(cr.StorageID); e != nil {
			o.Failure("storage_id.invalid", errInvalidStockStorage)
		}
	}

	if cr.ContainerID != "" {
		if cr.Container, e = validContainer(cr.ContainerID); e != nil {
			o.Failure("container_id.invalid", errInvalidContainer)
		}
	}

	return o
}

func (cr *finishRequest) Messages() map[string]string {
	return map[string]string{
		"container_id.required": errRequiredContainer,
	}
}

func (cr *finishRequest) Save() (e error) {
	u := cr.StockMovement
	u.FinishedAt = time.Now()
	u.Status = "finish"

	if e = u.Save("finished_at", "status"); e == nil {
		storage := takeStorage(cr.StockMovement.Destination, cr.Container, cr.Storage)

		if cr.StockMovement.IsPartial == 1 && cr.StockMovement.IsMerger != 1 {
			cr.StockMovement.NewUnit.SetStored(storage)
		} else if cr.StockMovement.IsPartial != 1 && cr.StockMovement.IsMerger == 1 {
			cr.StockMovement.Unit.SetVoid()
		} else if cr.StockMovement.IsPartial != 1 && cr.StockMovement.IsMerger != 1 {
			cr.StockMovement.Unit.SetStored(storage)
		}

		go event.Call("stockmovement::finished", cr.StockMovement)
	}

	return
}

func takeStorage(l *warehouse.Location, c *model2.Item, st *model.StockStorage) (storage *model.StockStorage) {
	if st != nil {
		return st
	}

	storage = new(model.StockStorage)
	storage.Location = l
	storage.Container = c
	storage.Save()

	return
}
