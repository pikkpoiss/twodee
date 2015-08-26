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
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/raster"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
)

type TextCache struct {
	Text     string
	Texture  *Texture
	fontface *FontFace
}

func NewTextCache(fontface *FontFace) *TextCache {
	return &TextCache{
		fontface: fontface,
	}
}

func (tc *TextCache) SetText(text string) (err error) {
	if text == tc.Text {
		return
	}
	var tex *Texture = tc.Texture
	tc.Texture, err = tc.fontface.GetText(text)
	if tex != nil {
		tex.Delete()
	}
	if err == nil {
		tc.Text = text
	}
	return
}

func (tc *TextCache) Clear() {
	tc.Delete()
	tc.Text = ""
}

func (tc *TextCache) Delete() {
	if tc.Texture != nil {
		tc.Texture.Delete()
		tc.Texture = nil
	}
}

type FontFace struct {
	font    *truetype.Font
	charw   float32
	charh   float32
	fg      color.Color
	bg      color.Color
	context *freetype.Context
}

func NewFontFace(path string, size float64, fg, bg color.Color) (fontface *FontFace, err error) {
	var (
		font      *truetype.Font
		fontbytes []byte
		bounds    truetype.Bounds
		context   *freetype.Context
		scale     float32
	)
	if fontbytes, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if font, err = freetype.ParseFont(fontbytes); err != nil {
		return
	}
	bounds = font.Bounds(1)
	context = freetype.NewContext()
	context.SetFont(font)
	context.SetFontSize(size)
	context.SetDPI(72)
	scale = float32(context.PointToFixed(size) >> 8)
	fontface = &FontFace{
		font:    font,
		charw:   scale * float32(bounds.XMax-bounds.XMin),
		charh:   scale * float32(bounds.YMax-bounds.YMin),
		fg:      fg,
		bg:      bg,
		context: context,
	}
	return
}

func (ff *FontFace) GetText(text string) (t *Texture, err error) {
	var (
		src       image.Image
		bg        image.Image
		dst       draw.Image
		shortened draw.Image
		pt        fixed.Point26_6
		w         int
		h         int
	)
	src = image.NewUniform(ff.fg)
	bg = image.NewUniform(ff.bg)
	w = int(float32(len(text)) * ff.charw)
	h = int(ff.charh)
	dst = image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(dst, dst.Bounds(), bg, image.ZP, draw.Src)
	ff.context.SetSrc(src)
	ff.context.SetDst(dst)
	ff.context.SetClip(dst.Bounds())
	pt = freetype.Pt(0, int(ff.charh))
	if pt, err = ff.context.DrawString(text, pt); err != nil {
		return
	}
	// if err = WritePNG("hello.png", dst); err != nil {
	// 	return
	// }
	shortened = image.NewRGBA(image.Rect(0, 0, int(pt.X/256), h))
	draw.Draw(shortened, shortened.Bounds(), dst, image.ZP, draw.Src)
	t, err = GetTexture(shortened, gl.NEAREST)
	return
}
