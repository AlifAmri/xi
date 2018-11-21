// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/gudang/api/src/auth"
	model2 "git.qasico.com/gudang/api/src/partnership/model"
	"git.qasico.com/gudang/api/src/receiving/model"
	"git.qasico.com/gudang/api/src/receiving/services"
	"git.qasico.com/gudang/api/src/user"

	"git.qasico.com/cuxs/validation"
)

type updateRequest struct {
	ID           int64   `json:"-" valid:"required"`
	PartnerID    string  `json:"partner_id" valid:"required"`
	DocumentCode string  `json:"document_code" valid:"required"`
	DocumentFile string  `json:"document_file"`
	Note         string  `json:"note"`
	Items        []*item `json:"items" valid:"required"`

	Session           *auth.SessionData   `json:"-"`
	Receiving         *model.Receiving    `json:"-"`
	Partner           *model2.Partnership `json:"-"`
	TotalQuantityPlan float64             `json:"-"`
}

func (ur *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if ur.Receiving, e = validReceiving(ur.ID, "active"); e != nil {
		o.Failure("id.invalid", errInvalidReceiving)
	}

	if ur.DocumentCode != "" && ur.Receiving != nil {
		if !validDocumentCode(ur.DocumentCode, ur.Receiving.ID) {
			o.Failure("document_code.unique", errUniqueCode)
		}
	}

	if ur.PartnerID != "" {
		if ur.Partner, e = validPartner(ur.PartnerID); e != nil {
			o.Failure("partner_id.invalid", errInvalidPartner)
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

func (ur *updateRequest) Save() (u *model.Receiving, e error) {
	ur.Receiving.DocumentCode = ur.DocumentCode
	ur.Receiving.DocumentFile = ur.DocumentFile
	ur.Receiving.Partner = ur.Partner
	ur.Receiving.Supervisor = ur.Session.User.(*user.User)
	ur.Receiving.Note = ur.Note
	ur.Receiving.TotalQuantityPlan = ur.TotalQuantityPlan

	if e = ur.Receiving.Save("document_code", "note", "partner_id", "document_file", "supervisor_id", "total_quantity_plan"); e == nil {
		orm.NewOrm().LoadRelated(ur.Receiving, "Documents", 0)

		for _, item := range ur.Items {
			item.Save(ur.Receiving)
		}

		// delete item yang tidak lagi di proses
		if ur.Receiving.Documents != nil {
			for _, document := range ur.Receiving.Documents {
				var used bool
				for _, item := range ur.Items {
					if item.ReceivingDocument != nil && document.ID == item.ReceivingDocument.ID {
						used = true
					}
				}

				if !used {
					document.Delete()
				}
			}
		}

		go services.CreateActual(ur.Receiving)
		go func() {
			o := orm.NewOrm()
			o.Raw("UPDATE stock_movement SET ref_code = ? "+
				"where ref_id = ? and type = 'putaway' and id > 0;", ur.Receiving.DocumentCode, ur.Receiving.ID).Exec()
			o.Raw("UPDATE incoming_vehicle ivh INNER JOIN receiving r ON r.vehicle_id = ivh.id "+
				"SET ivh.status = 'in_progress' WHERE r.id = ? AND ivh.status = ?", ur.ID, "in_queue").Exec()
		}()
	}

	return ur.Receiving, e
}
