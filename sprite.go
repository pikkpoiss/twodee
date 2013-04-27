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
)

type Sprite struct {
	Element
	system    *System
	texture   *Texture
	frame     int
	texture1  float64
	texture2  float64
	textureB  float64
	VelocityX float64
	VelocityY float64
	Type      int
	Collide   bool
	FlipX     bool
}

func (s *System) NewSprite(name string, x, y, w, h float64, t int) *Sprite {
	sprite := &Sprite{
		frame:   -1,
		system:  s,
		texture: s.Textures[name],
		textureB: 0.0,
		Type:    t,
		FlipX:   false,
		Collide: true,
	}
	sprite.SetBounds(Rect(x, y, x+w, y+h))
	sprite.SetFrame(0)
	return sprite
}

func (s *Sprite) SetTextureHeight(h float64) {
	s.textureB = 1 - (h / float64(s.texture.Height))
}

func (s *Sprite) SetVelocity(pt Point) {
	s.VelocityX = pt.X
	s.VelocityY = pt.Y
}

func (s *Sprite) Velocity() Point {
	return Pt(s.VelocityX, s.VelocityY)
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
	if frame == s.frame {
		return
	}
	s.frame = frame % len(s.texture.Frames)
	var (
		tex = s.texture.Frames[s.frame]
	)
	s.texture1 = float64(tex[0]) / float64(s.texture.Width)
	s.texture2 = float64(tex[1]) / float64(s.texture.Width)
}

func (s *Sprite) Draw() {
	var (
		b    = s.GlobalBounds()
		z    = s.Z()
		minx = b.Min.X
		maxx = b.Max.X
	)
	if s.FlipX {
		minx = b.Max.X
		maxx = b.Min.X
	}
	s.texture.Bind()
	gl.MatrixMode(gl.TEXTURE)
	gl.Begin(gl.QUADS)
	gl.TexCoord2d(s.texture1, s.textureB)
	gl.Vertex3d(minx, b.Min.Y, z)
	gl.TexCoord2d(s.texture2, s.textureB)
	gl.Vertex3d(maxx, b.Min.Y, z)
	gl.TexCoord2d(s.texture2, 1)
	gl.Vertex3d(maxx, b.Max.Y, z)
	gl.TexCoord2d(s.texture1, 1)
	gl.Vertex3d(minx, b.Max.Y, z)
	gl.End()
	gl.MatrixMode(gl.MODELVIEW)
	s.texture.Unbind()
	s.Element.Draw()
}

func (s *Sprite) Update() {
	s.MoveTo(Pt(s.X()+s.VelocityX, s.Y()+s.VelocityY))
}
