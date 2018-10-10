// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package preparation

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/auth"
	"git.qasico.com/gudang/api/src/delivery/model"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/user"
	"time"

	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID           int64   `json:"-" valid:"required"`
	PartnerID    string  `json:"partner_id" valid:"required"`
	DocumentCode string  `json:"document_code" valid:"required"`
	DocumentFile string  `json:"document_file"`
	ShippedDate  string  `json:"shipped_date" valid:"required"`
	Note         string  `json:"note"`
	Items        []*item `json:"items" valid:"required"`

	Session           *auth.SessionData   `json:"-"`
	Preparation       *model.Preparation  `json:"-"`
	Partner           *model2.Partnership `json:"-"`
	ShippedAt         time.Time           `json:"-"`
	TotalQuantityPlan float64             `json:"-"`
}

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.Preparation, e = validPreparation(ur.ID, "draft"); e != nil {
		o.Failure("id.invalid", errInvalidPreparation)
	}

	if ur.DocumentCode != "" && ur.Preparation != nil {
		if !validDocumentCode(ur.DocumentCode, ur.Preparation.ID) {
			o.Failure("document_code.unique", errUniqueCode)
		}
	}

	if ur.PartnerID != "" {
		if ur.Partner, e = validPartner(ur.PartnerID); e != nil {
			o.Failure("partner_id.invalid", errInvalidPartner)
		}
	}

	if ur.ShippedDate != "" {
		if ur.ShippedAt, e = time.Parse("2006-01-02", ur.ShippedDate); e != nil {
			o.Failure("shipped_date.invalid", errInvalidDate)
		}
	}

	if len(ur.Items) > 0 {
		for i, item := range ur.Items {
			item.Validate(i, o)
			ur.TotalQuantityPlan += item.Quantity
		}
	}

	return o
}

func (ur *updateRequest) Messages() map[string]string {
	return map[string]string{
		"partner_id.required":    errRequiredPartner,
		"document_code.required": errRequiredDocumentCode,
		"items.required":         errRequiredItems,
	}
}

func (ur *updateRequest) Save() (u *model.Preparation, e error) {
	ur.Preparation.DocumentCode = ur.DocumentCode
	ur.Preparation.DocumentFile = ur.DocumentFile
	ur.Preparation.Partner = ur.Partner
	ur.Preparation.Supervisor = ur.Session.User.(*user.User)
	ur.Preparation.Note = ur.Note
	ur.Preparation.TotalQuantityPlan = ur.TotalQuantityPlan
	ur.Preparation.ShippedAt = ur.ShippedAt

	if e = ur.Preparation.Save("document_code", "note", "partner_id", "document_file", "supervisor_id", "total_quantity_plan", "shipped_at"); e == nil {
		orm.NewOrm().LoadRelated(ur.Preparation, "Documents", 0)

		for _, item := range ur.Items {
			item.Save(ur.Preparation)
		}

		// delete item yang tidak lagi di proses
		if ur.Preparation.Documents != nil {
			for _, document := range ur.Preparation.Documents {
				var used bool
				for _, item := range ur.Items {
					if item.PreparationDocument != nil && document.ID == item.PreparationDocument.ID {
						used = true
					}
				}

				if !used {
					document.Delete()
				}
			}
		}
	}

	return ur.Preparation, e
}
