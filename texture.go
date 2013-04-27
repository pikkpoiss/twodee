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
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type Texture struct {
	texture gl.Texture
	Width   int
	Height  int
	Frames  [][]int
}

func LoadTexture(path string, smoothing int, framewidth int) (texture *Texture, err error) {
	var (
		img    image.Image
		bounds image.Rectangle
		gltex  gl.Texture
	)
	if img, err = LoadPNG(path); err != nil {
		return
	}
	img = GetPow2Image(img)
	bounds = img.Bounds()
	if gltex, err = GetGLTexture(img, smoothing); err != nil {
		return
	}
	texture = &Texture{
		texture: gltex,
		Width:   bounds.Dx(),
		Height:  bounds.Dy(),
		Frames:  make([][]int, 0),
	}
	frames := bounds.Dx() / framewidth
	for i := 0; i < frames; i++ {
		texture.Frames = append(texture.Frames, []int{
			i * framewidth,
			(i + 1) * framewidth,
		})
	}
	return
}

func LoadVarWidthTexture(path string, smoothing int) (texture *Texture, err error) {
	var (
		img   image.Image
		trim  *image.NRGBA
		gltex gl.Texture
	)
	if img, err = LoadPNG(path); err != nil {
		return
	}
	var (
		bounds     = img.Bounds()
		trimbounds = image.Rect(0, 0, bounds.Dx(), bounds.Dy()-1)
		trimpoint  = image.Pt(0, 1)
	)
	trim = image.NewNRGBA(trimbounds)
	draw.Draw(trim, trimbounds, img, trimpoint, draw.Src)
	if gltex, err = GetGLTexture(trim, smoothing); err != nil {
		return
	}
	texture = &Texture{
		texture: gltex,
		Width:   trimbounds.Dx(),
		Height:  trimbounds.Dy(),
		Frames:  make([][]int, 0),
	}
	var (
		aprime uint32 = 0
		pair          = make([]int, 2)
		x             = 0
	)
	for ; x < bounds.Dx(); x++ {
		_, _, _, a := img.At(x, 0).RGBA()
		if aprime == 0 && a > 0 {
			pair[0] = x
		} else if aprime > 0 && a == 0 {
			pair[1] = x
			texture.Frames = append(texture.Frames, pair)
			pair = make([]int, 2)
		}
		aprime = a
	}
	if pair[0] != 0 {
		pair[1] = x
		texture.Frames = append(texture.Frames, pair)
	}
	return
}

func (t *Texture) Bind() {
	t.texture.Bind(gl.TEXTURE_2D)
}

func (t *Texture) Unbind() {
	t.texture.Unbind(gl.TEXTURE_2D)
}

func (t *Texture) Dispose() {
	t.texture.Delete()
}

func getPow2(i int) int {
	p2 := 1
	for p2 < i {
		p2 = p2 << 1
	}
	return p2
}

func GetPow2Image(img image.Image) image.Image {
	var (
		b   = img.Bounds()
		p2w = getPow2(b.Max.X)
		p2h = getPow2(b.Max.Y)
	)
	if p2w == b.Max.X && p2h == b.Max.Y {
		return img
	}
	out := image.NewRGBA(image.Rect(0, 0, p2w, p2h))
	draw.Draw(out, b, img, image.ZP, draw.Src)
	return out
}

func LoadPNG(path string) (img image.Image, err error) {
	var file *os.File
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()
	img, err = png.Decode(file)
	return
}

func GetGLTexture(img image.Image, smoothing int) (gltexture gl.Texture, err error) {
	var data *bytes.Buffer
	if data, err = EncodeTGA("texture", img); err != nil {
		return
	}
	gltexture = gl.GenTexture()
	gltexture.Bind(gl.TEXTURE_2D)
	if !glfw.LoadMemoryTexture2D(data.Bytes(), 0) {
		err = fmt.Errorf("Failed to load texture")
		return
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, smoothing)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, smoothing)
	return
}
