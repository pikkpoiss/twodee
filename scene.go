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
)

type Node interface {
	AddChild(node Node)
	Draw()
}

type Element struct {
	Children []Node
}

func (e *Element) AddChild(node Node) {
	e.Children = append(e.Children, node)
}

func (e *Element) Draw() {
	for _, child := range e.Children {
		child.Draw()
	}
}

type Scene struct {
	Element
}

type Sprite struct {
	Element
	system  *System
	texture *Texture
	X       float32
	Y       float32
	Width   int
	Height  int
	Frames  int
	Frame   int
}

func (s *System) NewSprite(name string, x float32, y float32, w int, h int, frames int) *Sprite {
	if frames <= 0 {
		frames = 1
	}
	return &Sprite{
		system:  s,
		texture: s.Textures[name],
		X:       x,
		Y:       y,
		Width:   w,
		Height:  h,
		Frame:   0,
		Frames:  frames,
	}
}

func (s *Sprite) Draw() {
	s.Element.Draw()
	var (
		x2         float32 = s.X + float32(s.Width)
		y2         float32 = s.Y + float32(s.Height)
		frame      int     = s.Frame % s.Frames
		framestep  float32 = 1.0 / float32(s.Frames)
		framestart float32 = framestep * float32(frame)
		framestop  float32 = framestep * float32(frame+1)
	)
	s.texture.Bind()
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(framestop, 1)
	gl.Vertex2f(s.X, s.Y)
	gl.TexCoord2f(framestop, 0)
	gl.Vertex2f(x2, s.Y)
	gl.TexCoord2f(framestart, 0)
	gl.Vertex2f(x2, y2)
	gl.TexCoord2f(framestart, 1)
	gl.Vertex2f(s.X, y2)
	gl.End()
	s.texture.Unbind()
}
