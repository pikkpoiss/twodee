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

const (
	Step60Hz = time.Duration(16666) * time.Microsecond
	Step30Hz = Step60Hz * 2
	Step20Hz = time.Duration(50000) * time.Microsecond
	Step15Hz = Step30Hz * 2
	Step10Hz = Step20Hz * 2
	Step5Hz = Step10Hz * 2
)

type Animation struct {
	accumulated time.Duration
	Current     int
	FrameLength time.Duration
	Sequence    []int
	callback    AnimationCallback
}

func NewAnimation(length time.Duration, frames []int) *Animation {
	return &Animation{
		FrameLength: length,
		Sequence:    frames,
		Current:     frames[0],
	}
}

func (a *Animation) Update(elapsed time.Duration) {
	a.accumulated += elapsed
	index := int(a.accumulated/a.FrameLength) % len(a.Sequence)
	a.Current = a.Sequence[index]
	if a.callback != nil && index == len(a.Sequence)-1 {
		a.callback()
		a.callback = nil
	}
}

func (a *Animation) OffsetFrame(offset int) int {
	index := int(a.accumulated/a.FrameLength) % len(a.Sequence)
	return a.Sequence[(index + offset) % len(a.Sequence)]
}

func (a *Animation) SetCallback(callback AnimationCallback) {
	a.callback = callback
}

type AnimationCallback func()

func (a *Animation) SetSequence(seq []int) {
	a.Sequence = seq
	a.Current = a.Sequence[0]
	a.accumulated = time.Duration(0)
}
