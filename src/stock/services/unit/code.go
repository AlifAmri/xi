// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package unit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"git.qasico.com/cuxs/orm"
)

const codeFormat = "B%08d"
const codeFormatChild = "%s-%d"

// NewCode membuat code unit stock baru
func NewCode() string {
	return fmt.Sprintf(codeFormat, lastCode())
}

// NewChildCode membuat code unit stock baru
func NewChildCode(parent string) string {
	return fmt.Sprintf(codeFormatChild, parent, lastChild(fmt.Sprintf("%s-", parent)))
}

func lastCode() int {
	// ambil id terakhir dari code
	o := orm.NewOrm()

	var lastCode string
	var lastNumber int
	if e := o.Raw("select code from stock_unit where code like ? and code not like ? order by code desc", "B%", "%-%").QueryRow(&lastCode); e == nil {
		number := regexp.MustCompile(`[\d]+$`).FindString(lastCode)
		lastNumber, _ = strconv.Atoi(number)
	}

	return lastNumber + 1
}

func lastChild(code string) int {
	o := orm.NewOrm()

	var lastCode string
	var lastNumber int
	if e := o.Raw("select code from stock_unit where code like ? order by code desc", "%"+code+"%").QueryRow(&lastCode); e == nil {
		number := strings.Replace(lastCode, code, "", 1)
		lastNumber, _ = strconv.Atoi(number)
	}

	return lastNumber + 1
}
