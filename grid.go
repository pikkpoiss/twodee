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
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
)

type GridItem interface {
	Passable() bool
	Opaque() bool
}

type Grid struct {
	Width     int32
	Height    int32
	BlockSize float32
	points    []GridItem
}

func NewGrid(w, h, blocksize int32) *Grid {
	return &Grid{
		points:    make([]GridItem, w*h),
		Width:     w,
		Height:    h,
		BlockSize: float32(blocksize),
	}
}

func (g *Grid) Index(x, y int32) int32 {
	if x < 0 || y < 0 {
		return -1
	}
	return g.Width*(g.Height-y-1) + x
}

func (g *Grid) Get(x, y int32) GridItem {
	return g.GetIndex(g.Index(x, y))
}

func (g *Grid) GetIndex(index int32) GridItem {
	if index < 0 || index > g.Width*g.Height {
		return nil
	}
	return g.points[index]
}

func (g *Grid) Set(x, y int32, val GridItem) {
	g.SetIndex(g.Index(x, y), val)
}

func (g *Grid) SetIndex(index int32, val GridItem) {
	if index < 0 || index > g.Width*g.Height {
		return
	}
	g.points[index] = val
}

func (g *Grid) GetImage(fg, bg color.Color) *image.NRGBA {
	var (
		img  = image.NewNRGBA(image.Rect(0, 0, int(g.Width), int(g.Height)))
		item GridItem
	)
	for x := 0; x < int(g.Width); x++ {
		for y := 0; y < int(g.Height); y++ {
			item = g.Get(int32(x), int32(y))
			if item != nil && item.Passable() {
				img.Set(x, y, fg)
			} else {
				img.Set(x, y, bg)
			}
		}
	}
	return img
}

func (g *Grid) squareCollides(bounds mgl32.Vec4, x, y float32) bool {
	// Bounds are {minx, miny, maxx, maxy}
	// Sizex, sizey are the number of coordinate units a grid entry occupies.
	var (
		size  = g.BlockSize
		fudge = float32(0.001) // Prevents item from sticking to wall when we round its coordinates.
		minx  = int32((bounds[0] + x) / size)
		miny  = int32((bounds[1] + y) / size)
		maxx  = int32((bounds[2] + x - fudge) / size)
		maxy  = int32((bounds[3] + y - fudge) / size)
		i     int32
		j     int32
		item  GridItem
	)
	for i = minx; i <= maxx; i++ {
		for j = miny; j <= maxy; j++ {
			item = g.Get(i, j)
			if item != nil && item.Passable() {
				return true
			}
		}
	}
	return false
}

func (g *Grid) FixMove(bounds mgl32.Vec4, move mgl32.Vec2) (out mgl32.Vec2) {
	out = move
	if g.squareCollides(bounds, out[0], 0.0) {
		out[0] = g.GridAligned(bounds[0]) - bounds[0]
	}
	if g.squareCollides(bounds, out[0], out[1]) {
		out[1] = g.GridAligned(bounds[1]) - bounds[1]
	}
	return
}

func (g *Grid) GridAligned(x float32) float32 {
	return g.BlockSize * float32(int32((x/g.BlockSize)+0.5))
}

func (g *Grid) GridPosition(v float32) int32 {
	return int32(v / g.BlockSize)
}

func (g *Grid) InversePosition(i int32) float32 {
	return float32(i)*g.BlockSize + g.BlockSize/2.0
}

func (g *Grid) CanSee(from, to mgl32.Vec2) bool {
	var (
		size  = g.BlockSize
		minx  = int32(from[0] / size)
		maxx  = int32(to[0] / size)
		miny  = int32(from[1] / size)
		maxy  = int32(to[1] / size)
		slope = float32(maxy-miny) / float32(maxx-minx)
		c     = float32(miny) - (slope * float32(minx))
		x     int32
		y     int32
		item  GridItem
	)
	for x = minx; x <= maxx; x++ {
		y = int32(slope*float32(x) + c)
		item = g.Get(x, y)
		if item != nil && item.Opaque() {
			// Something blocks the way
			return false
		}
	}
	for y = miny; y <= maxy; y++ {
		x = int32((float32(y) - c) / slope)
		item = g.Get(x, y)
		if item != nil && item.Opaque() {
			// Something blocks the way
			return false
		}
	}
	return true
}
