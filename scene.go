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
	parent Node
	X float32
	Y float32
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
	gl.TexCoord2f(framestop, 1)
	gl.Vertex2f(x1, y1)
	gl.TexCoord2f(framestop, 0)
	gl.Vertex2f(x2, y1)
	gl.TexCoord2f(framestart, 0)
	gl.Vertex2f(x2, y2)
	gl.TexCoord2f(framestart, 1)
	gl.Vertex2f(x1, y2)
	gl.End()
	s.texture.Unbind()
}

type Environment struct {
	Element
}

func (s *System) LoadEnvironment(name string, path string) (env *Environment, err error) {
	var (
		file   *os.File
		img    image.Image
		bounds image.Rectangle
		colors map[uint32]int
		sprite *Sprite
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()
	if img, err = png.Decode(file); err != nil {
		return
	}
	colors = make(map[uint32]int, 0)
	bounds = img.Bounds()
	env = &Environment{}
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			index := (r << 8 + g) << 8 + b
			if _, exists := colors[index]; exists == false {
				colors[index] = len(colors)
			}
			sprite = s.NewSprite(name, float32(x * 32), float32(y * 32), 32, 32, 2)
			sprite.Frame = colors[index]
			env.AddChild(sprite)
		}
	}
	return
}
