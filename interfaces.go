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

type Spatial interface {
	X() float64
	Y() float64
	Height() float64
	Width() float64
	MoveTo(p Point)
	Bounds() Rectangle
	SetWidth(w float64)
	SetHeight(h float64)
	SetVelocity(p Point)
	Velocity() Point
}

type Visible interface {
	Draw()
}

type Changing interface {
	Update()
}

type SpatialChanging interface {
	Spatial
	Changing
}

type SpatialVisible interface {
	Spatial
	Visible
}

type SpatialVisibleChanging interface {
	Spatial
	Visible
	Changing
}

type MapLoader interface {
	Create(tileset string, index int, x, y, w, h float64)
	Loaded(bounds Rectangle, properties map[string]string)
}

type Node interface {
	AddChild(node Node)
	RemoveChild(node Node)
	GetAllChildren() []Node
	Parent() Node
	SetParent(Node)
	Draw()
	SetBounds(Rectangle)
	GlobalBounds() Rectangle
	Bounds() Rectangle
	RelativeBounds(Node) Rectangle
	Width() float64
	Height() float64
	SetWidth(float64)
	SetHeight(float64)
	SetZ(float64)
	X() float64
	Y() float64
	Z() float64
}
