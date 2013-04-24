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

type Positionable interface {
	X() float64
	Y() float64
	SetX(x float64)
	SetY(y float64)
	Bounds() Rectangle
	SetWidth(w float64)
	SetHeight(h float64)
}

type Updateable interface {
	Update()
}

type SpriteFactory interface {
	Create(tileset string, index int, x, y, w, h float64) *Sprite
}
