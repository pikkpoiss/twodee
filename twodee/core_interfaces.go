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
	"time"
)

type Animating interface {
	Update(elapsed time.Duration) (done bool)
	SetCallback(callback AnimationCallback)
	HasCallback() bool
	Reset()
}

type Entity interface {
	Pos() Point
	Bounds() Rectangle
	Frame() int
	Rotation() float32
	MoveTo(Point)
	MoveToCoords(x, y float32)
	Update(elapsed time.Duration)
}

type Event interface{}

type GETyper interface {
	GEType() GameEventType
}

type Layer interface {
	Render()
	Update(elapsed time.Duration)
	Delete()
	HandleEvent(evt Event) bool
	Reset() error
}

type MenuItem interface {
	Highlighted() bool
	Label() string
	setHighlighted(val bool)
	Parent() MenuItem
	setParent(item MenuItem)
	Active() bool
	setActive(val bool)
}

type TexturedTile interface {
	ScaledBounds(ratio float32) (x, y, w, h float32)
	ScaledTextureBounds(rx float32, ry float32) (x, y, w, h float32)
}
