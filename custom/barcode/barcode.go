// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package barcode

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"git.qasico.com/cuxs/env"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code39"
	"github.com/disintegration/imaging"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// MakeBarcode creating an image barcode from the code
func MakeBarcode(code string) (r string, e error) {
	var bcode barcode.Barcode

	if bcode, e = code39.Encode(code, false, false); e != nil {
		return
	}

	if bcode, e = barcode.Scale(bcode, 290, 40); e != nil {
		return
	}

	newImg := imaging.New(290, 54, color.NRGBA{255, 255, 255, 255})
	newImg = imaging.Paste(newImg, bcode, image.Pt(0, 0))
	newImg = imaging.Paste(newImg, addLabel(code), image.Pt(0, 42))

	dir := env.GetString("BARCODE_DIRECTORY", ".")
	path := env.GetString("BARCODE_PATH", "")

	if e = draw2dimg.SaveToPngFile(fmt.Sprintf("%s/%s.png", dir, code), newImg); e == nil {
		r = fmt.Sprintf("%s/%s.png", path, code)
	}

	return
}

func addLabel(code string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 250, 10))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	point := fixed.Point26_6{X: fixed.Int26_6(110 * 64), Y: fixed.Int26_6(10 * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString(code)

	return img
}
