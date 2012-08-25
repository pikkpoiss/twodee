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
}

type Element struct {
	Children []Node
	parent   Node
	X        float32
	Y        float32
}

func (e *Element) AddChild(node Node) {
	node.SetParent(e)
	e.Children = append(e.Children, node)
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
	if e.parent == nil {
		return e.X
	}
	return e.parent.GlobalX() + e.X
}

func (e *Element) GlobalY() float32 {
	if e.parent == nil {
		return e.Y
	}
	return e.parent.GlobalY() + e.Y
}

type Scene struct {
	Element
}

type Sprite struct {
	Element
	system  *System
	texture *Texture
	Width   int
	Height  int
	Frames  int
	Frame   int
}

func (s *System) NewSprite(name string, x float32, y float32, w int, h int, frames int) *Sprite {
	if frames <= 0 {
		frames = 1
	}
	sprite := &Sprite{
		system:  s,
		texture: s.Textures[name],
		Width:   w,
		Height:  h,
		Frame:   0,
		Frames:  frames,
	}
	sprite.X = x
	sprite.Y = y
	return sprite
}

func (s *Sprite) Draw() {
	s.Element.Draw()
	var (
		x1         float32 = s.GlobalX()
		y1         float32 = s.GlobalY()
		x2         float32 = x1 + float32(s.Width)
		y2         float32 = y1 + float32(s.Height)
		frame      int     = s.Frame % s.Frames
		framestep  float32 = 1.0 / float32(s.Frames)
		framestart float32 = framestep * float32(frame)
		framestop  float32 = framestep * float32(frame+1)
	)
	s.texture.Bind()
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(framestart, 1)
	gl.Vertex2f(x1, y1)
	gl.TexCoord2f(framestop, 1)
	gl.Vertex2f(x2, y1)
	gl.TexCoord2f(framestop, 0)
	gl.Vertex2f(x2, y2)
	gl.TexCoord2f(framestart, 0)
	gl.Vertex2f(x1, y2)
	gl.End()
	s.texture.Unbind()
}

type EnvOpts struct {
	Blocks      []*EnvBlock
	TextureName string
	MapPath     string
	BlockWidth  int
	BlockHeight int
	Frames      int
}

type EnvBlockLoadedHandler func(sprite *Sprite, block *EnvBlock)

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
			)
			if block, exists = colors[index]; exists == false {
				// Unrecognized colors just get a pass
				continue
			}
			sprite = s.NewSprite(
				opts.TextureName,
				float32(x*opts.BlockWidth),
				float32(y*opts.BlockHeight),
				opts.BlockWidth,
				opts.BlockHeight,
				opts.Frames)
			sprite.Frame = block.FrameIndex
			env.AddChild(sprite)
			if block.Handler != nil {
				block.Handler(sprite, block)
			}
		}
	}
	return
}
