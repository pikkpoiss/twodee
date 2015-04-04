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
	"bytes"
	"fmt"
	"github.com/go-gl/gl/v3.2-core/gl"
	"image"
)

type Texture struct {
	Texture        uint32
	Width          int
	Height         int
	OriginalWidth  int
	OriginalHeight int
}

func LoadTexture(path string, smoothing int) (texture *Texture, err error) {
	var (
		img image.Image
	)
	if img, err = loadPNG(path); err != nil {
		return
	}
	return GetTexture(img, smoothing)
}

func GetTexture(img image.Image, smoothing int) (texture *Texture, err error) {
	var bounds = img.Bounds()
	return GetTruncatedTexture(img, smoothing, bounds.Dx(), bounds.Dy())
}

func GetTruncatedTexture(img image.Image, smoothing int, w, h int) (texture *Texture, err error) {
	var (
		bounds image.Rectangle
		gltex  uint32
	)
	img = getPow2Image(img)
	bounds = img.Bounds()
	if gltex, err = GetGLTexture(img, smoothing); err != nil {
		return
	}
	texture = &Texture{
		Texture:        gltex,
		Width:          bounds.Dx(),
		Height:         bounds.Dy(),
		OriginalWidth:  w,
		OriginalHeight: h,
	}
	return
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.Texture)
}

func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t *Texture) Delete() {
	if t.Texture != 0 {
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.DeleteTextures(1, &t.Texture)
		t.Texture = 0
	}
}

func GetGLTexture(img image.Image, smoothing int) (t uint32, err error) {
	var (
		data   *bytes.Buffer
		bounds image.Rectangle
		width  int
		height int
	)
	if data, err = imageBytes(img); err != nil {
		return
	}
	bounds = img.Bounds()
	width = bounds.Max.X - bounds.Min.X
	height = bounds.Max.Y - bounds.Min.Y
	gl.GenTextures(1, &t)
	if e := gl.GetError(); e != 0 {
		fmt.Printf("ggt1 ERROR: %s\n", e)
	}
	gl.BindTexture(gl.TEXTURE_2D, t)
	if e := gl.GetError(); e != 0 {
		fmt.Printf("ggt2 ERROR: %s\n", e)
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(smoothing))
	if e := gl.GetError(); e != 0 {
		fmt.Printf("ggt3 ERROR: %s\n", e)
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(smoothing))
	if e := gl.GetError(); e != 0 {
		fmt.Printf("ggt4 ERROR: %s\n", e)
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_INT_8_8_8_8, gl.Ptr(data.Bytes()))
	if e := gl.GetError(); e != 0 {
		fmt.Printf("ggt5 ERROR: %s\n", e)
	}
	gl.GenerateMipmap(gl.TEXTURE_2D)
	if e := gl.GetError(); e != 0 {
		fmt.Printf("ggt6 ERROR: %s\n", e)
	}
	gl.BindTexture(gl.TEXTURE_2D, 0)
	if e := gl.GetError(); e != 0 {
		fmt.Printf("ggt7 ERROR: %s\n", e)
	}
	return
}
