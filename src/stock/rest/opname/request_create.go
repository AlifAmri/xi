// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package opname

import (
	"fmt"
	"time"

	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/stock/model"
	"git.qasico.com/gudang/api/src/user"
	"git.qasico.com/gudang/api/src/warehouse"

	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	LocationID string        `json:"location_id" valid:"required"`
	Note       string        `json:"note"`
	Items      []*opnameItem `json:"items" valid:"required"`

	Location *warehouse.Location `json:"-"`
	Session  *auth.SessionData   `json:"-"`
}

func (cr *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if cr.LocationID != "" {
		if cr.Location, e = validLocation(cr.LocationID); e != nil {
			o.Failure("location_id.invalid", errInvalidLocation)
		}
	}

	if len(cr.Items) > 0 {
		contNum := make(map[int8]string)
		contcollector := make(map[int8]bool)
		for i, item := range cr.Items {
			item.Validate(i, o)
			if item.IsVoid == int8(0) {
				contcollector[item.ContainerNum] = true
			}
			if item.ContainerNum > 0 && item.ContainerID != "" {
				if val, ok := contNum[item.ContainerNum]; ok {
					if item.ContainerID != val {
						o.Failure(fmt.Sprintf("items.%d.container_num.invalid", i), errRequiredContainerNum)
						o.Failure(fmt.Sprintf("items.%d.container_id.invalid", i), errInvalidContainer)
					}
				} else {
					contNum[item.ContainerNum] = item.ContainerID
				}
			}
		}
		// cek max pallet dan container
		if (len(contcollector) + countMovement(cr.Location.ID)) > cr.Location.StorageCapacity {
			o.Failure("location_id.invalid", "pallet di lokasi ini sudah maksimum")
		}
	}

	return o
}

func (cr *createRequest) Messages() map[string]string {
	return map[string]string{
		"location_id.required": errRequiredLocation,
		"items.required":       errRequiredItems,
	}
}

func (cr *createRequest) Save() (u *model.StockOpname, e error) {
	u = &model.StockOpname{
		Location:  cr.Location,
		Type:      "opname",
		Status:    "active",
		Note:      cr.Note,
		CreatedBy: cr.Session.User.(*user.User),
		CreatedAt: time.Now(),
	}

	if e = u.Save(); e == nil {
		for _, item := range cr.Items {
			item.Save(u, cr.Session.User.(*user.User))
		}
	}

	return
}
