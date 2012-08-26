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
	"github.com/banthar/gl"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Node interface {
	AddChild(node Node)
	Parent() Node
	SetParent(Node)
	Draw()
	GlobalX() float32
	GlobalY() float32
	GlobalZ() float32
	IsOffsetParent() bool
}

type Element struct {
	Children []Node
	parent   Node
	X        float32
	Y        float32
	Z        float32
}

func (e *Element) IsOffsetParent() bool {
	return false
}

func (e *Element) AddChild(node Node) {
	node.SetParent(e)
	e.Children = append(e.Children, node)
}

func (e *Element) Clear() {
	e.Children = make([]Node, 0)
}

func (e *Element) SetParent(node Node) {
	e.parent = node
}

func (e *Element) Parent() Node {
	return e.parent
}

func (e *Element) Draw() {
	for _, child := range e.Children {
		child.Draw()
	}
}

func (e *Element) GlobalX() float32 {
	if e.parent == nil || e.parent.IsOffsetParent() {
		return e.X
	}
	return e.parent.GlobalX() + e.X
}

func (e *Element) GlobalY() float32 {
	if e.parent == nil || e.parent.IsOffsetParent() {
		return e.Y
	}
	return e.parent.GlobalY() + e.Y
}

func (e *Element) GlobalZ() float32 {
	if e.parent == nil || e.parent.IsOffsetParent() {
		return e.Z
	}
	return e.parent.GlobalZ() + e.Z
}

type Scene struct {
	Element
}

type Sprite struct {
	Element
	system    *System
	texture   *Texture
	Width     int
	Height    int
	frame     int
	texture1  float64
	texture2  float64
	VelocityX float32
	VelocityY float32
	Type      int
}

func (s *System) NewSprite(name string, x float32, y float32, w int, h int, t int) *Sprite {
	sprite := &Sprite{
		system:  s,
		texture: s.Textures[name],
		Width:   w,
		Height:  h,
		Type:    t,
	}
	sprite.X = x
	sprite.Y = y
	sprite.SetFrame(0)
	return sprite
}

func inside(point float32, start float32, end float32) bool {
	return start < point && point < end
}

func (s *Sprite) TestMove(dx float32, dy float32, sprite *Sprite) bool {
	var (
		x11 float32 = s.GlobalX() + dx
		x12 float32 = x11 + float32(s.Width)
		x21 float32 = sprite.GlobalX()
		x22 float32 = x21 + float32(sprite.Width)
	)
	inX := (inside(x11, x21, x22) || inside(x12, x21, x22) || inside(x21, x11, x12) || inside(x22, x11, x12))
	if !inX && x11 != x21 && x12 != x22 {
		return true
	}
	var (
		y11 float32 = s.GlobalY() + dy
		y12 float32 = y11 + float32(s.Height)
		y21 float32 = sprite.GlobalY()
		y22 float32 = y21 + float32(sprite.Height)
	)
	inY := (inside(y11, y21, y22) || inside(y12, y21, y22) || inside(y21, y11, y12) || inside(y22, y11, y12))
	if !inY {
		return true
	}
	return false
}

func (s *Sprite) CollidesWith(sprite *Sprite) bool {
	return !s.TestMove(0, 0, sprite)
}

func (s *Sprite) SetFrame(frame int) {
	s.frame = frame % len(s.texture.Frames)
	s.texture1 = float64(s.texture.Frames[s.frame][0]) / float64(s.texture.Width)
	s.texture2 = float64(s.texture.Frames[s.frame][1]) / float64(s.texture.Width)
}

func (s *Sprite) Draw() {
	var (
		x1 float32 = s.GlobalX()
		y1 float32 = s.GlobalY()
		z1 float32 = s.GlobalZ()
		x2 float32 = x1 + float32(s.Width)
		y2 float32 = y1 + float32(s.Height)
	)
	s.texture.Bind()
	gl.MatrixMode(gl.TEXTURE)
	//gl.Scalef(1.0/64.0, 1.0/64.0, 1.0);
	gl.Begin(gl.QUADS)
	gl.TexCoord2d(s.texture1, 1)
	gl.Vertex3f(x1, y1, z1)
	gl.TexCoord2d(s.texture2, 1)
	gl.Vertex3f(x2, y1, z1)
	gl.TexCoord2d(s.texture2, 0)
	gl.Vertex3f(x2, y2, z1)
	gl.TexCoord2d(s.texture1, 0)
	gl.Vertex3f(x1, y2, z1)
	gl.End()
	gl.MatrixMode(gl.MODELVIEW)
	s.texture.Unbind()
	s.Element.Draw()
}

type Text struct {
	Element
	system  *System
	texture *Texture
	Width   int
	Height  int
	text    string
	ratio   int
}

func (s *System) NewText(name string, x float32, y float32, r int, text string) *Text {
	t := &Text{
		system:  s,
		ratio:   r,
		texture: s.Textures[name],
	}
	t.X = x
	t.Y = y
	t.Height = t.texture.Height * t.ratio
	t.SetText(text)
	return t
}

func (t *Text) SetText(text string) {
	t.Clear()
	var x int = 0
	for _, c := range text {
		frame := (int(c) - int(' ')) % len(t.texture.Frames)
		width := t.ratio * (t.texture.Frames[frame][1] - t.texture.Frames[frame][0])
		sprite := &Sprite{
			system:  t.system,
			texture: t.texture,
			Width:   width,
			Height:  t.Height,
		}
		sprite.SetFrame(frame)
		sprite.X = float32(x)
		sprite.Y = 0
		x += width + (1 * t.ratio)
		t.AddChild(sprite)
	}
	t.Width = x
}

type EnvOpts struct {
	Blocks      []*EnvBlock
	TextureName string
	MapPath     string
	BlockWidth  int
	BlockHeight int
}

type EnvBlockLoadedHandler func(env *Env, block *EnvBlock, sprite *Sprite, x float32, y float32)

type EnvBlock struct {
	Type       int
	Color      color.Color
	FrameIndex int
	Handler    EnvBlockLoadedHandler
}

type Env struct {
	Element
	Width  int
	Height int
}

func (s *System) LoadEnv(opts EnvOpts) (env *Env, err error) {
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
	env = &Env{
		Width:  opts.BlockWidth * bounds.Dx(),
		Height: opts.BlockHeight * bounds.Dy(),
	}
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			index := GetIndex(img.At(x, y))
			var (
				block  *EnvBlock
				exists bool
				gX     = float32(x * opts.BlockWidth)
				gY     = float32(y * opts.BlockHeight)
			)
			if block, exists = colors[index]; exists == false {
				// Unrecognized colors just get a pass
				continue
			}
			// Pass -1 to not render anything (important parts)
			if block.FrameIndex != -1 {
				sprite = s.NewSprite(
					opts.TextureName,
					gX,
					gY,
					opts.BlockWidth,
					opts.BlockHeight,
					block.Type)
				sprite.SetFrame(block.FrameIndex)
				env.AddChild(sprite)
			}
			if block.Handler != nil {
				block.Handler(env, block, sprite, gX, gY)
			}
		}
	}
	return
}

func (e *Env) IsOffsetParent() bool {
	return true
}

