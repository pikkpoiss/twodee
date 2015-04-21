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
	"math"
	"time"
)

const (
	Step60Hz = time.Duration(16666) * time.Microsecond
	Step30Hz = Step60Hz * 2
	Step20Hz = time.Duration(50000) * time.Microsecond
	Step15Hz = Step30Hz * 2
	Step10Hz = Step20Hz * 2
	Step5Hz  = Step10Hz * 2
)

type AnimationCallback func()

type Animation struct {
	elapsed  time.Duration
	callback AnimationCallback
}

func NewAnimation() *Animation {
	return &Animation{}
}

func (a *Animation) Update(elapsed time.Duration) (done bool) {
	a.elapsed += elapsed
	return false
}

func (a *Animation) Elapsed() time.Duration {
	return a.elapsed
}

func (a *Animation) Reset() {
	a.elapsed = time.Duration(0)
}

func (a *Animation) SetCallback(callback AnimationCallback) {
	//if a.callback != nil {
	//	a.callback()
	//}
	a.callback = callback
}

func (a *Animation) HasCallback() bool {
	return a.callback != nil
}

func (a *Animation) Callback() {
	if a.HasCallback() {
		a.callback()
	}
}

type FrameAnimation struct {
	*Animation
	FrameLength time.Duration
	Sequence    []int
	Current     int
}

func NewFrameAnimation(length time.Duration, frames []int) *FrameAnimation {
	return &FrameAnimation{
		Animation:   NewAnimation(),
		FrameLength: length,
		Sequence:    frames,
		Current:     frames[0],
	}
}

func (a *FrameAnimation) Update(elapsed time.Duration) (done bool) {
	a.Animation.Update(elapsed)
	index := int(a.Elapsed()/a.FrameLength) % len(a.Sequence)
	a.Current = a.Sequence[index]
	done = false
	if a.HasCallback() && index == len(a.Sequence)-1 {
		a.Callback()
		a.SetCallback(nil)
		done = true
	}
	return
}

func (a *FrameAnimation) OffsetFrame(offset int) int {
	index := int(a.Elapsed()/a.FrameLength) % len(a.Sequence)
	return a.Sequence[(index+offset)%len(a.Sequence)]
}

func (a *FrameAnimation) SetSequence(seq []int) {
	a.Sequence = seq
	a.Current = a.Sequence[0]
	a.Animation.Reset()
}

type ContinuousFunc func(elapsed time.Duration) float32

type ContinuousAnimation struct {
	*Animation
	function ContinuousFunc
}

func NewContinuousAnimation(f ContinuousFunc) *ContinuousAnimation {
	return &ContinuousAnimation{
		Animation: NewAnimation(),
		function:  f,
	}
}

func (a *ContinuousAnimation) Value() float32 {
	return a.function(a.Elapsed())
}

func SineDecayFunc(duration time.Duration, amplitude, frequency, decay float32, callback AnimationCallback) ContinuousFunc {
	var interval = float64(frequency * 2.0 * math.Pi)
	return func(elapsed time.Duration) float32 {
		if elapsed > duration {
			if callback != nil {
				callback()
			}
			return 0.0
		}
		decayAmount := 1.0 - float32(elapsed)/float32(duration)*decay
		return float32(math.Sin(elapsed.Seconds()*interval/duration.Seconds())) * amplitude * decayAmount
	}
}
