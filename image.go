// Copyright 2013 Arne Roomann-Kurrik
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
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
)

func EncodeTGA(name string, img image.Image) (buf *bytes.Buffer, err error) {
	var (
		bounds image.Rectangle = img.Bounds()
		ident  []byte          = []byte(name)
		width  []byte          = make([]byte, 2)
		height []byte          = make([]byte, 2)
		nrgba  *image.NRGBA
		data   []byte
	)
	binary.LittleEndian.PutUint16(width, uint16(bounds.Dx()))
	binary.LittleEndian.PutUint16(height, uint16(bounds.Dy()))

	// See http://paulbourke.net/dataformats/tga/
	buf = &bytes.Buffer{}
	buf.WriteByte(byte(len(ident)))
	buf.WriteByte(0)
	buf.WriteByte(2) // uRGBI
	buf.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0})
	buf.Write([]byte(width))
	buf.Write([]byte(height))
	buf.WriteByte(32) // Bits per pixel
	buf.WriteByte(8)
	if buf.Len() != 18 {
		err = fmt.Errorf("TGA header is not 18 bytes: %v", buf.Len())
		return
	}

	nrgba = image.NewNRGBA(bounds)
	draw.Draw(nrgba, bounds, img, bounds.Min, draw.Src)
	buf.Write(ident)
	data = make([]byte, bounds.Dx()*bounds.Dy()*4)
	var (
		lineLength int = bounds.Dx() * 4
		destOffset int = len(data) - lineLength
	)
	for srcOffset := 0; srcOffset < len(nrgba.Pix); {
		var (
			dest   = data[destOffset : destOffset+lineLength]
			source = nrgba.Pix[srcOffset : srcOffset+nrgba.Stride]
		)
		copy(dest, source)
		destOffset -= lineLength
		srcOffset += nrgba.Stride
	}
	for x := 0; x < len(data); {
		buf.WriteByte(data[x+2])
		buf.WriteByte(data[x+1])
		buf.WriteByte(data[x+0])
		buf.WriteByte(data[x+3])
		x += 4
	}
	return
}
