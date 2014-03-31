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

import ()

type Point struct {
	X float32
	Y float32
}

func Pt(x, y float32) Point {
	return Point{x, y}
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
