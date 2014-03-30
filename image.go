// Copyright 2014 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package twodee

import (
	"bufio"
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func loadPNG(path string) (img image.Image, err error) {
	var file *os.File
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()
	img, err = png.Decode(file)
	return
}

func imageBytes(img image.Image) (buf *bytes.Buffer, err error) {
	var (
		bounds image.Rectangle = img.Bounds()
		rgba   *image.RGBA
		data   []byte
	)
	buf = &bytes.Buffer{}
	rgba = image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	data = make([]byte, len(rgba.Pix))
	var (
		destOffset int = len(data) - rgba.Stride
	)
	for srcOffset := 0; srcOffset < len(rgba.Pix); {
		var (
			dest   = data[destOffset : destOffset+rgba.Stride]
			source = rgba.Pix[srcOffset : srcOffset+rgba.Stride]
		)
		copy(dest, source)
		destOffset -= rgba.Stride
		srcOffset += rgba.Stride
	}
	for x := 0; x < len(data); {
		buf.WriteByte(data[x+3])
		buf.WriteByte(data[x+2])
		buf.WriteByte(data[x+1])
		buf.WriteByte(data[x+0])
		x += 4
	}
	return
}

func pow2(i int) int {
	p2 := 1
	for p2 < i {
		p2 = p2 << 1
	}
	return p2
}

func getPow2Image(img image.Image) image.Image {
	var (
		bounds = img.Bounds()
		width  = pow2(bounds.Max.X)
		height = pow2(bounds.Max.Y)
	)
	if width == bounds.Max.X && height == bounds.Max.Y {
		return img
	}
	out := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(out, bounds, img, image.ZP, draw.Src)
	return out
}

func WritePNG(path string, img image.Image) (err error) {
	var (
		f *os.File
		b *bufio.Writer
	)
	if f, err = os.Create(path); err != nil {
		return
	}
	defer f.Close()
	b = bufio.NewWriter(f)
	if err = png.Encode(b, img); err != nil {
		return
	}
	err = b.Flush()
	return
}
