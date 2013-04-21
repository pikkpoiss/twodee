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
	View Rectangle
}

func NewCamera(x float64, y float64, w float64, h float64) *Camera {
	return &Camera{
		View: Rect(x, y, x+w, y+h),
	}
}

func (c *Camera) MatchRatio(width int, height int) {
	ratio := float64(height) / float64(width)
	c.View.Max.Y = c.View.Max.X * ratio
}

func (c *Camera) SetProjection() {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(c.View.Min.X, c.View.Max.X, c.View.Max.Y, c.View.Min.Y, -1, 1)
}
