// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package barcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeBarcode(t *testing.T) {
	s, e := MakeBarcode("001390057")

	assert.NoError(t, e)
	assert.NotNil(t, s)
}
