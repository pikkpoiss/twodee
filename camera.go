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

type Camera struct {
	view    Rectangle
	focus   Point
	width   float64
	height  float64
	zoom    float64
	limits  Rectangle
	limited bool
}

func NewCamera(x float64, y float64, w float64, h float64) (c *Camera) {
	c = &Camera{
		width:   w,
		height:  h,
		focus:   Pt(x+w/2.0, y+h/2.0),
		zoom:    0,
		limited: false,
	}
	c.calcView()
	return
}

func (c *Camera) calcView() bool {
	var (
		ratio = c.height / c.width
		hw    = c.width / 2.0
		hh    = hw * ratio
		zw    = hw * c.zoom
		zh    = zw * ratio
		view  = c.view
		focus = c.focus
	)
	c.view.Min.X = c.focus.X - hw - zw
	c.view.Min.Y = c.focus.Y - hh - zh
	c.view.Max.X = c.focus.X + hw + zw
	c.view.Max.Y = c.focus.Y + hh + zh

	if c.limited {
		if !c.view.In(c.limits) {
			if c.view.Max.X > c.limits.Max.X {
				c.focus.X -= c.view.Max.X - c.limits.Max.X
			}
			if c.view.Min.X < c.limits.Min.X {
				c.focus.X -= c.view.Min.X - c.limits.Min.X
			}
			if c.view.Max.Y > c.limits.Max.Y {
				c.focus.Y -= c.view.Max.Y - c.limits.Max.Y
			}
			if c.view.Min.Y < c.limits.Min.Y {
				c.focus.Y -= c.view.Min.Y - c.limits.Min.Y
			}
			if c.view.In(c.limits) {
				c.view = view
				c.focus = focus
				return false
			}
			//TODO: Terrible hack for keeping the camera from freezing
			//TODO: Fix zooming with real maths
			c.limited = false
			c.calcView()
			c.limited = true
			return true
		}
		return true
	}
	return true
}

func (c *Camera) MatchRatio(width int, height int) {
	h := c.height
	ratio := float64(height) / float64(width)
	c.height = c.width * ratio
	if !c.calcView() {
		c.height = h
	}
}

func (c *Camera) Bottom(y float64) {
	var (
		dy = y - c.view.Min.Y
	)
	c.focus.Y += dy
	if !c.calcView() {
		c.focus.Y -= dy
	}
}

func (c *Camera) Pan(x float64, y float64) {
	c.focus.X += x
	c.focus.Y += y
	if !c.calcView() {
		c.focus.X -= x
		c.focus.Y -= y
	}
}

func (c *Camera) Focus() Point {
	return c.focus
}

func (c *Camera) SetFocus(x float64, y float64) {
	var (
		ox = c.focus.X
		oy = c.focus.Y
	)
	c.focus.X = x
	c.focus.Y = y
	if !c.calcView() {
		c.focus.X = ox
		c.focus.Y = oy
	}
}

func (c *Camera) Zoom(incr float64) {
	z := c.zoom
	c.zoom += incr
	if c.zoom < -0.9 {
		c.zoom = z
		return
	}
	if !c.calcView() {
		c.zoom = z
	}
}

func (c *Camera) Bounds() Rectangle {
	return c.view
}

func (c *Camera) SetLimits(limits Rectangle) {
	c.limits = limits
	c.limited = true
}

func (c *Camera) ResolveScreenCoords(x, y, w, h int) (gx, gy float64) {
	gx = c.view.Min.X + (float64(x) / float64(w)) * (c.view.Max.X - c.view.Min.X)
	gy = c.view.Min.Y + (float64(h - y) / float64(h)) * (c.view.Max.Y - c.view.Min.Y)
	return
}

func (c *Camera) SetProjection() {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(c.view.Min.X, c.view.Max.X, c.view.Min.Y, c.view.Max.Y, 1, -1)
}
