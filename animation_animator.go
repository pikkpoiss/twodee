// Copyright 2015 Twodee Authors
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
	"math"
	"time"
)

type AnimatorCallback func()

type Animator interface {
	SetCallback(Callback AnimatorCallback)
	IsDone() bool
	Update(elapsed time.Duration) time.Duration
	Reset()
	Delete()
}

type BoundedAnimation struct {
	Elapsed  time.Duration
	Duration time.Duration
	Callback AnimatorCallback
}

func (a *BoundedAnimation) SetCallback(Callback AnimatorCallback) {
	a.Callback = Callback
}

func (a *BoundedAnimation) IsDone() bool {
	return a.Elapsed >= a.Duration
}

func (a *BoundedAnimation) Update(elapsed time.Duration) time.Duration {
	a.Elapsed += elapsed
	if a.IsDone() {
		if a.Callback != nil {
			a.Callback()
		}
		return a.Elapsed - a.Duration
	}
	return 0
}

func (a *BoundedAnimation) Reset() {
	a.Elapsed = 0
}

func (a *BoundedAnimation) Delete() {
}

type GroupedAnimation struct {
	animators []Animator
	Callback  AnimatorCallback
}

func (a *GroupedAnimation) SetCallback(Callback AnimatorCallback) {
	a.Callback = Callback
}

func (a *GroupedAnimation) IsDone() bool {
	var done = true
	for _, animator := range a.animators {
		if !animator.IsDone() {
			done = false
		}
	}
	return done
}

func (a *GroupedAnimation) Update(elapsed time.Duration) time.Duration {
	var (
		total     time.Duration
		remainder time.Duration
		done      = true
	)
	for _, animator := range a.animators {
		remainder = animator.Update(elapsed)
		if !animator.IsDone() {
			done = false
		}
		if remainder != 0 && (total == 0 || remainder < total) {
			total = remainder // Take the smallest nonzero remainder.
		}
	}
	if done {
		if a.Callback != nil {
			a.Callback()
		}
		return total
	}
	return 0
}

func (a *GroupedAnimation) Reset() {
	for _, animator := range a.animators {
		animator.Reset()
	}
}

func (a *GroupedAnimation) Delete() {
	for _, animator := range a.animators {
		animator.Delete()
	}
	a.animators = []Animator{}
}

type ChainedAnimation struct {
	animators []Animator
	loop      bool
	index     int
	Callback  AnimatorCallback
}

func (a *ChainedAnimation) SetCallback(Callback AnimatorCallback) {
	a.Callback = Callback
}

func (a *ChainedAnimation) IsDone() bool {
	var count = len(a.animators)
	return !a.loop && count > 0 && a.animators[count-1].IsDone()
}

func (a *ChainedAnimation) Update(elapsed time.Duration) time.Duration {
	var count = len(a.animators)
	if count > a.index {
		for elapsed > 0 && !a.animators[a.index].IsDone() {
			elapsed = a.animators[a.index].Update(elapsed)
			if a.animators[a.index].IsDone() {
				if a.loop {
					a.animators[a.index].Reset()
				}
				a.index = (a.index + 1) % count
				if !a.loop && a.index == 0 && a.Callback != nil {
					a.Callback()
					break
				}
			}
		}
	}
	return elapsed
}

func (a *ChainedAnimation) Reset() {
	a.index = 0
	for _, animator := range a.animators {
		animator.Reset()
	}
}

func (a *ChainedAnimation) Delete() {
	for _, animator := range a.animators {
		animator.Delete()
	}
	a.animators = []Animator{}
}

type LinearAnimation struct {
	BoundedAnimation
	target *float32
	from   float32
	to     float32
}

func NewLinearAnimation(target *float32, from, to float32, duration time.Duration) *LinearAnimation {
	return &LinearAnimation{
		BoundedAnimation{
			0,
			duration,
			nil,
		},
		target,
		from,
		to,
	}
}

func (a *LinearAnimation) Update(elapsed time.Duration) (remainder time.Duration) {
	remainder = a.BoundedAnimation.Update(elapsed)
	var (
		denom = float64(a.BoundedAnimation.Duration)
		numer = math.Min(float64(a.BoundedAnimation.Elapsed), denom)
		pct   = float32(numer / denom)
		value = pct*(a.to-a.from) + a.from
	)
	*a.target = value
	return
}

type EaseOutAnimation struct {
	BoundedAnimation
	target *float32
	from   float32
	to     float32
}

func NewEaseOutAnimation(target *float32, from, to float32, duration time.Duration) *EaseOutAnimation {
	return &EaseOutAnimation{
		BoundedAnimation{
			0,
			duration,
			nil,
		},
		target,
		from,
		to,
	}
}

func (a *EaseOutAnimation) Update(elapsed time.Duration) (remainder time.Duration) {
	remainder = a.BoundedAnimation.Update(elapsed)
	var (
		denom = float64(a.BoundedAnimation.Duration)
		numer = math.Min(float64(a.BoundedAnimation.Elapsed), denom)
		pct   = float32(numer / denom)
		c     = a.to - a.from
		value = -c*pct*(pct-2) + a.from
	)
	*a.target = value
	return
}
