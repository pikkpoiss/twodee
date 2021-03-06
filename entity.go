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

type BaseEntity struct {
	pos      Point
	halfW    float32
	halfH    float32
	rotation float32
	frame    int
}

func NewBaseEntity(x, y, w, h, r float32, frame int) *BaseEntity {
	return &BaseEntity{
		pos:      Pt(x, y),
		halfW:    w / 2.0,
		halfH:    h / 2.0,
		rotation: r,
		frame:    frame,
	}
}

func (e *BaseEntity) Bounds() Rectangle {
	return Rect(
		e.pos.X()-e.halfW,
		e.pos.Y()-e.halfH,
		e.pos.X()+e.halfW,
		e.pos.Y()+e.halfH,
	)
}

func (e *BaseEntity) Pos() Point {
	return e.pos
}

func (e *BaseEntity) MoveTo(pt Point) {
	e.pos = pt
}

func (e *BaseEntity) MoveToCoords(x, y float32) {
	e.pos = Pt(x, y)
}

func (e *BaseEntity) Frame() int {
	return e.frame
}

func (e *BaseEntity) Rotation() float32 {
	return e.rotation
}

func (e *BaseEntity) Update(elapsed time.Duration) {
}

type AnimatingEntity struct {
	animation *FrameAnimation
	*BaseEntity
}

func NewAnimatingEntity(x, y, w, h, r float32, l time.Duration, f []int) *AnimatingEntity {
	return &AnimatingEntity{
		animation:  NewFrameAnimation(l, f),
		BaseEntity: NewBaseEntity(x, y, w, h, r, 0),
	}
}

func (e *AnimatingEntity) SetFrames(f []int) {
	e.animation.SetSequence(f)
}

func (e *AnimatingEntity) SetCallback(callback AnimationCallback) {
	e.animation.SetCallback(callback)
}

func (e *AnimatingEntity) Update(elapsed time.Duration) {
	e.animation.Update(elapsed)
}

func (e *AnimatingEntity) Frame() int {
	return e.animation.Current
}

func (e *AnimatingEntity) OffsetFrame(offset int) int {
	return e.animation.OffsetFrame(offset)
}
