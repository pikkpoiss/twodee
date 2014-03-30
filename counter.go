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

type Counter struct {
	Count int64
	Total float64
	Last  float64
	Avg   float32
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Incr() {
	now := glfw.GetTime()
	if int32(c.Last) != int32(now) {
		c.Total += math.Remainder(now-c.Last, 1)
		c.Count += 1
		avg := c.Total / float64(c.Count)
		c.Avg = float32(avg * 1000)
		c.Total = math.Remainder(now, 1)
		c.Count = 1
	} else {
		c.Total += (now - c.Last)
		c.Count += 1
	}
	c.Last = now
}
