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
	SetBounds(Rectangle)
	GlobalBounds() Rectangle
	Bounds() Rectangle
	RelativeBounds(Node) Rectangle
	Width() float32
	Height() float32
	SetWidth(float32)
	SetHeight(float32)
	SetZ(float32)
	X() float32
	Y() float32
	Z() float32
}

type Element struct {
	Children []Node
	parent   Node
	z        float32
	bounds   Rectangle
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

func (e *Element) SetBounds(r Rectangle) {
	e.bounds = r
}

func (e *Element) GlobalBounds() Rectangle {
	if e.parent == nil {
		return e.Bounds()
	}
	return e.Bounds().Add(e.parent.GlobalBounds().Min)
}

func (e *Element) Bounds() Rectangle {
	return e.bounds
}

func (e *Element) RelativeBounds(n Node) Rectangle {
	return e.GlobalBounds().Sub(n.GlobalBounds().Min)
}

func (e *Element) Move(p Point) {
	e.bounds = e.bounds.Add(p)
}

func (e *Element) MoveTo(p Point) {
	var (
		x = p.X + e.bounds.Max.X - e.bounds.Min.X
		y = p.Y + e.bounds.Max.Y - e.bounds.Min.Y
	)
	e.bounds = Rectangle{Min: p, Max: Pt(x, y)}
}

func (e *Element) Width() float32 {
	return e.bounds.Max.X - e.bounds.Min.X
}

func (e *Element) SetWidth(w float32) {
	e.bounds.Max.X = e.bounds.Min.X + w
}

func (e *Element) Height() float32 {
	return e.bounds.Max.Y - e.bounds.Min.Y
}

func (e *Element) SetHeight(h float32) {
	e.bounds.Max.Y = e.bounds.Min.Y + h
}

func (e *Element) X() float32 {
	return e.bounds.Min.X
}

func (e *Element) Y() float32 {
	return e.bounds.Min.Y
}

func (e *Element) Z() float32 {
	if e.parent == nil {
		return e.z
	}
	return e.parent.Z() + e.z
}

func (e *Element) SetZ(z float32) {
	e.z = z
}

type Scene struct {
	Element
}

type Sprite struct {
	Element
	system    *System
	texture   *Texture
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
		Type:    t,
	}
	sprite.SetBounds(Rect(x, y, x+float32(w), y+float32(h)))
	sprite.SetFrame(0)
	return sprite
}

func (s *Sprite) TestMove(dx float32, dy float32, r *Sprite) bool {
	var (
		sb = s.GlobalBounds()
		rb = r.GlobalBounds()
		p  = Pt(dx, dy)
	)
	if sb.Add(p).Overlaps(rb) {
		return false
	}
	return true
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
		b = s.GlobalBounds()
		z = s.Z()
	)
	s.texture.Bind()
	gl.MatrixMode(gl.TEXTURE)
	gl.Begin(gl.QUADS)
	gl.TexCoord2d(s.texture1, 1)
	gl.Vertex3f(b.Min.X, b.Min.Y, z)
	gl.TexCoord2d(s.texture2, 1)
	gl.Vertex3f(b.Max.X, b.Min.Y, z)
	gl.TexCoord2d(s.texture2, 0)
	gl.Vertex3f(b.Max.X, b.Max.Y, z)
	gl.TexCoord2d(s.texture1, 0)
	gl.Vertex3f(b.Min.X, b.Max.Y, z)
	gl.End()
	gl.MatrixMode(gl.MODELVIEW)
	s.texture.Unbind()
	s.Element.Draw()
}

type Text struct {
	Element
	system  *System
	texture *Texture
	text    string
	ratio   int
}

func (s *System) NewText(name string, x float32, y float32, r int, text string) *Text {
	t := &Text{
		system:  s,
		ratio:   r,
		texture: s.Textures[name],
	}
	t.SetBounds(Rect(x, y, x, y+float32(t.texture.Height*t.ratio)))
	t.SetText(text)
	return t
}

func (t *Text) SetText(text string) {
	t.Clear()
	var x float32 = 0
	for _, c := range text {
		frame := (int(c) - int(' ')) % len(t.texture.Frames)
		width := t.ratio * (t.texture.Frames[frame][1] - t.texture.Frames[frame][0])
		sprite := &Sprite{
			system:  t.system,
			texture: t.texture,
		}
		sprite.SetBounds(Rect(x, 0, x+float32(width), t.Height()))
		sprite.SetFrame(frame)
		x += float32(width + (1 * t.ratio))
		t.AddChild(sprite)
	}
	t.SetWidth(float32(x))
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
	env = &Env{}
	env.SetBounds(Rect(0, 0, float32(opts.BlockWidth*bounds.Dx()), float32(opts.BlockHeight*bounds.Dy())))
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
