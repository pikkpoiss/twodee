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
	"github.com/go-gl/gl"
	"fmt"
)

type Sprite struct {
	Element
	system    *System
	texture   *Texture
	frame     int
	texture1  float64
	texture2  float64
	VelocityX float64
	VelocityY float64
	Type      int
	Ratio     float64
	Collide   bool
}

func (s *System) NewSprite(name string, x float64, y float64, w int, h int, t int) *Sprite {
	fmt.Printf("Textures: %v\n", s.Textures)
	sprite := &Sprite{
		system:  s,
		texture: s.Textures[name],
		Type:    t,
		//TODO: Figure out texture scaling in a better way
		Ratio:   float64(h) / float64(s.Textures[name].Height),
		Collide: true,
	}
	sprite.SetBounds(Rect(x, y, x+float64(w), y+float64(h)))
	sprite.SetFrame(0)
	return sprite
}

func (s *Sprite) TestMove(dx float64, dy float64, r *Sprite) bool {
	var (
		pad = float64(0.01)
		sb  = s.GlobalBounds()
		rb  = r.GlobalBounds()
		p   = Pt(dx, dy)
	)
	sb.Min.X += pad
	sb.Min.Y += pad
	sb.Max.X -= pad
	sb.Max.Y -= pad
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
	var (
		tex   = s.texture.Frames[s.frame]
		width = tex[1] - tex[0]
	)
	s.texture1 = float64(tex[0]) / float64(s.texture.Width)
	s.texture2 = float64(tex[1]) / float64(s.texture.Width)
	if s.Ratio != 0 {
		s.SetWidth(float64(width) * s.Ratio)
	}
}

func (s *Sprite) Draw() {
	var (
		b = s.GlobalBounds()
		z = s.Z()
	)
	s.texture.Bind()
	gl.MatrixMode(gl.TEXTURE)
	gl.Begin(gl.QUADS)
	gl.TexCoord2d(s.texture1, 0)
	gl.Vertex3d(b.Min.X, b.Min.Y, z)
	gl.TexCoord2d(s.texture2, 0)
	gl.Vertex3d(b.Max.X, b.Min.Y, z)
	gl.TexCoord2d(s.texture2, 1)
	gl.Vertex3d(b.Max.X, b.Max.Y, z)
	gl.TexCoord2d(s.texture1, 1)
	gl.Vertex3d(b.Min.X, b.Max.Y, z)
	gl.End()
	gl.MatrixMode(gl.MODELVIEW)
	s.texture.Unbind()
	s.Element.Draw()
}
