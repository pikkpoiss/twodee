// Copyright 2014 Arne Roomann-Kurrik
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
	"image"
	"image/color"
)

type Grid struct {
	Width  int32
	Height int32
	points []bool
}

func NewGrid(w, h int32) *Grid {
	return &Grid{
		points: make([]bool, w*h),
		Width:  w,
		Height: h,
	}
}

func (g *Grid) Index(x, y int32) int32 {
	return g.Width*y + x
}

func (g *Grid) Get(x, y int32) bool {
	return g.GetIndex(g.Index(x, y))
}

func (g *Grid) GetIndex(index int32) bool {
	if index < 0 || index > g.Width*g.Height {
		return false
	}
	return g.points[index]
}

func (g *Grid) Set(x, y int32, val bool) {
	g.SetIndex(g.Index(x, y), val)
}

func (g *Grid) SetIndex(index int32, val bool) {
	if index < 0 || index > g.Width*g.Height {
		return
	}
	g.points[index] = val
}

func (g *Grid) GetImage(fg, bg color.Color) *image.NRGBA {
	var img = image.NewNRGBA(image.Rect(0, 0, int(g.Width), int(g.Height)))
	for x := 0; x < int(g.Width); x++ {
		for y := 0; y < int(g.Height); y++ {
			if g.Get(int32(x), int32(y)) {
				img.Set(x, y, fg)
			} else {
				img.Set(x, y, bg)
			}
		}
	}
	return img
}
