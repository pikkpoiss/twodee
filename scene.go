// Copyright 2012 Arne Roomann-Kurrik
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

type ByDepth []Node

func (s ByDepth) Len() int {
	return len(s)
}

func (s ByDepth) Less(i int, j int) bool {
	return s[i].Z() < s[j].Z()
}

func (s ByDepth) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

type Scene struct {
	*Camera
	*Font
	Static  []SpatialVisible
	Dynamic []SpatialVisibleChanging
}

func (s *Scene) SetBounds(rect Rectangle) {
	s.Camera.SetLimits(rect)
}

func (s *Scene) Draw() {
	for _, e := range s.Static {
		e.Draw()
	}
	for _, e := range s.Dynamic {
		e.Draw()
	}
}

func (s *Scene) Update(u Updater) {
	for _, e := range s.Dynamic {
		u.Update(e)
	}
}
