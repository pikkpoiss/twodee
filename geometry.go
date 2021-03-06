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
)

type Point struct {
	mgl32.Vec2
}

func Pt(x, y float32) Point {
	return Point{mgl32.Vec2{x, y}}
}

func (p Point) Scale(a float32) Point {
	return Point{p.Vec2.Mul(a)}
}

func (p Point) Add(pt Point) Point {
	return Point{p.Vec2.Add(pt.Vec2)}
}

func (p Point) Sub(pt Point) Point {
	return Point{p.Vec2.Sub(pt.Vec2)}
}

func (p Point) DistanceTo(pt Point) float32 {
	return p.Sub(pt).Len()
}

type Rectangle struct {
	Min Point
	Max Point
}

func Rect(x1, y1, x2, y2 float32) Rectangle {
	return Rectangle{
		Min: Pt(x1, y1),
		Max: Pt(x2, y2),
	}
}

func (r Rectangle) Midpoint() Point {
	return Pt((r.Max.X()+r.Min.X())/2.0, (r.Max.Y()+r.Min.Y())/2.0)
}

func (r Rectangle) Overlaps(s Rectangle) bool {
	return s.Min.X() < r.Max.X() && s.Max.X() > r.Min.X() &&
		s.Min.Y() < r.Max.Y() && s.Max.Y() > r.Min.Y()
}

func (r Rectangle) ContainsPoint(a Point) bool {
	return r.Min.X() <= a.X() && r.Max.X() >= a.X() &&
		r.Min.Y() <= a.Y() && r.Max.Y() >= a.Y()
}

// Returns true if r is intersection by the line a, b.
func (r Rectangle) IntersectedBy(a, b Point) bool {
	if a.X() < r.Min.X() && b.X() < r.Min.X() {
		return false
	} else if a.X() > r.Max.X() && b.X() > r.Max.X() {
		return false
	} else if a.Y() < r.Min.Y() && b.Y() < r.Min.Y() {
		return false
	} else if a.Y() > r.Max.Y() && b.Y() > r.Max.Y() {
		return false
	} else {
		// The line is neither totally to the left, right, above, or below
		// the rectangle. There may be a collision.
		corners := []Point{
			Pt(r.Min.X(), r.Min.Y()),
			Pt(r.Min.Y(), r.Max.Y()),
			Pt(r.Max.X(), r.Min.Y()),
			Pt(r.Max.X(), r.Max.Y()),
		}
		eq := GetVectorDeterminantEquation(a, b)
		lastEvalSide := eq(corners[0]) > 0
		for _, corner := range corners[1:] {
			side := eq(corner) > 0
			if side != lastEvalSide {
				return true
			}
		}
	}
	return false
}

func GetVectorDeterminantEquation(a, b Point) func(Point) float32 {
	return func(p Point) float32 {
		return (p.X()-a.X())*(b.Y()-a.Y()) - (p.Y()-a.Y())*(b.X()-a.X())
	}
}
