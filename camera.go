// Copyright 2015 Arne Roomann-Kurrik
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
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	WorldBounds  Rectangle
	ScreenBounds Rectangle
	Projection   mgl32.Mat4
	Inverse      mgl32.Mat4
}

func NewCamera(world Rectangle, screen Rectangle) (c *Camera, err error) {
	c = &Camera{}
	c.SetScreenBounds(screen)
	err = c.SetWorldBounds(world)
	return
}

func (c *Camera) SetScreenBounds(bounds Rectangle) {
	c.ScreenBounds = bounds
}

func (c *Camera) SetWorldBounds(bounds Rectangle) (err error) {
	c.WorldBounds = bounds
	c.Projection = mgl32.Ortho(
		bounds.Min.X(),
		bounds.Max.X(),
		bounds.Min.Y(),
		bounds.Max.Y(),
		1,
		0)
	c.Inverse, err = GetInverseMatrix(c.Projection)
	return
}

func (c *Camera) ScreenToWorldCoords(x, y float32) (wx, wy float32) {
	// http://stackoverflow.com/questions/7692988/
	var (
		halfw = c.ScreenBounds.Max.X() / 2.0
		halfh = c.ScreenBounds.Max.Y() / 2.0
		xpct  = (x - halfw) / halfw
		ypct  = (halfh - y) / halfh
	)
	return Unproject(c.Inverse, xpct, ypct)
}

func (c *Camera) WorldToScreenCoords(x, y float32) (sx, sy float32) {
	var pctx, pcty = Project(c.Projection, x, y)
	var halfw = c.ScreenBounds.Max.X() / 2.0
	var halfh = c.ScreenBounds.Max.Y() / 2.0
	sx = pctx*halfw + halfw
	sy = pcty*halfh + halfh
	return
}
