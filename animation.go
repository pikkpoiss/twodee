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

import ()

type Animation struct {
	Frames  []int
	Repeat  int
	current int
	counter int
}

func Anim(frames []int, repeat int) *Animation {
	return &Animation{
		Frames:  frames,
		Repeat:  repeat,
		current: 0,
		counter: 0,
	}
}

func (a *Animation) Len() int {
	return len(a.Frames)
}

func (a *Animation) Next() int {
	a.counter += 1
	if a.counter > a.Repeat {
		a.current = (a.current + 1) % a.Len()
		a.counter = 0
	}
	return a.Frames[a.current]
}
