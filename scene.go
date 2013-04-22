// Copyright 2012 Arne Roomann-Kurrik
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
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sort"
)

func Round(a float64) float64 {
	if a > 0 {
		a += 0.5
	} else {
		a -= 0.5
	}
	return float64(math.Floor(float64(a)))
}

type Node interface {
	AddChild(node Node)
	RemoveChild(node Node)
	GetAllChildren() []Node
	Parent() Node
	SetParent(Node)
	Draw()
	SetBounds(Rectangle)
	GlobalBounds() Rectangle
	Bounds() Rectangle
	RelativeBounds(Node) Rectangle
	Width() float64
	Height() float64
	SetWidth(float64)
	SetHeight(float64)
	SetZ(float64)
	X() float64
	Y() float64
	Z() float64
}

type ByDepth []Node

func (s ByDepth) Len() int {
	return len(s)
}

func (s ByDepth) Less(i int, j int) bool {
	return s[i].Z() < s[j].Z()
}

func (s ByDepth) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

type Scene struct {
	Element
	*Camera
	*Font
}

func (s *Scene) Draw() {
	l := s.GetAllChildren()
	sort.Sort(ByDepth(l))
	for _, c := range l {
		c.Draw()
	}
}

type Text struct {
	Element
	system  *System
	texture *Texture
	text    string
	ratio   int
}

func (s *System) NewText(name string, x float64, y float64, r int, text string) *Text {
	t := &Text{
		system:  s,
		ratio:   r,
		texture: s.Textures[name],
	}
	t.SetBounds(Rect(x, y, x, y+float64(t.texture.Height*t.ratio)))
	t.SetText(text)
	return t
}

func (t *Text) SetText(text string) {
	t.Clear()
	var x float64 = 0
	for _, c := range text {
		frame := (int(c) - int(' ')) % len(t.texture.Frames)
		width := t.ratio * (t.texture.Frames[frame][1] - t.texture.Frames[frame][0])
		sprite := &Sprite{
			system:  t.system,
			texture: t.texture,
		}
		sprite.SetBounds(Rect(x, 0, x+float64(width), t.Height()))
		sprite.SetFrame(frame)
		x += float64(width + (1 * t.ratio))
		t.AddChild(sprite)
	}
	t.SetWidth(float64(x))
}

type EnvOpts struct {
	Blocks      []*EnvBlock
	TextureName string
	MapPath     string
	BlockWidth  int
	BlockHeight int
}

type EnvBlockLoadedHandler func(block *EnvBlock, sprite *Sprite, x float64, y float64)

type EnvBlock struct {
	Type       int
	Color      color.Color
	FrameIndex int
	Handler    EnvBlockLoadedHandler
}

type Env struct {
	Element
}

func (e *Env) Load(system *System, opts EnvOpts) (err error) {
	var (
		file   *os.File
		img    image.Image
		bounds image.Rectangle
		colors map[uint32]*EnvBlock
		sprite *Sprite
	)
	GetIndex := func(c color.Color) uint32 {
		r, g, b, a := c.RGBA()
		return ((r<<8+g)<<8+b)<<8 + a
	}
	if file, err = os.Open(opts.MapPath); err != nil {
		return
	}
	defer file.Close()
	if img, err = png.Decode(file); err != nil {
		return
	}
	colors = make(map[uint32]*EnvBlock, 0)
	for _, block := range opts.Blocks {
		colors[GetIndex(block.Color)] = block
	}
	bounds = img.Bounds()
	e.SetBounds(Rect(0, 0, float64(opts.BlockWidth*bounds.Dx()), float64(opts.BlockHeight*bounds.Dy())))
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			index := GetIndex(img.At(x, y))
			var (
				block  *EnvBlock
				exists bool
				gX     = float64(x * opts.BlockWidth)
				gY     = float64(y * opts.BlockHeight)
			)
			if block, exists = colors[index]; exists == false {
				// Unrecognized colors just get a pass
				continue
			}
			// Pass -1 to not render anything (important parts)
			if block.FrameIndex != -1 {
				sprite = system.NewSprite(
					opts.TextureName,
					gX,
					gY,
					opts.BlockWidth,
					opts.BlockHeight,
					block.Type)
				sprite.SetFrame(block.FrameIndex)
				e.AddChild(sprite)
			}
			if block.Handler != nil {
				block.Handler(block, sprite, gX, gY)
			}
		}
	}
	return
}
