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
	lastNum := checkUniqueCodeParent(codeFormat, lastCode())
	return fmt.Sprintf(codeFormat, lastNum)
}

// NewChildCode membuat code unit stock baru
func NewChildCode(parent string) string {
	lastNum := checkUniqueCodeChild(codeFormatChild, parent, lastChild(fmt.Sprintf("%s-", parent)))
	return fmt.Sprintf(codeFormatChild, parent, lastNum)
}

func lastCode() int {
	// ambil id terakhir dari code
	o := orm.NewOrm()

	var lastCode string
	var lastNumber int
	if e := o.Raw("select code from stock_unit where code like ? and code not like ? order by id desc", "B%", "%-%").QueryRow(&lastCode); e == nil {
		number := regexp.MustCompile(`[\d]+$`).FindString(lastCode)
		lastNumber, _ = strconv.Atoi(number)
	}

	return lastNumber + 1
}

func lastChild(code string) int {
	o := orm.NewOrm()

	var lastCode string
	var lastNumber int
	if e := o.Raw("select code from stock_unit where code like ? order by id desc", "%"+code+"%").QueryRow(&lastCode); e == nil {
		number := strings.Replace(lastCode, code, "", 1)
		lastNumber, _ = strconv.Atoi(number)
	}

	return lastNumber + 1
}

func checkUniqueCodeParent(prefix string, last int) int {
	var result int64
	code := fmt.Sprintf(prefix, last)
	o := orm.NewOrm()
	o.Raw("SELECT id FROM stock_unit WHERE code = ?", code).QueryRow(&result)
	if result != int64(0) {
		last = checkUniqueCodeParent(prefix, last+1)
	}
	return last
}

func checkUniqueCodeChild(prefix string, parent string, last int) int {
	var result int64
	code := fmt.Sprintf(prefix, parent, last)
	o := orm.NewOrm()
	o.Raw("SELECT id FROM stock_unit WHERE code = ?", code).QueryRow(&result)
	if result != int64(0) {
		last = checkUniqueCodeChild(prefix, parent, last+1)
	}
	return last
}
