// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

var (
	errRequiredName    = "nama aplikasi harus diisi"
	errRequiredCompany = "nama perusahaan harus diisi"
)

type attribute struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}
