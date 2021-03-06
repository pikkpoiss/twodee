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
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

func GetInverseMatrix(m mgl32.Mat4) (out mgl32.Mat4, err error) {
	var (
		empty = mgl32.Mat4{}
	)
	if out = m.Inv(); out == empty {
		err = fmt.Errorf("Matrix %v not invertible", m)
		return
	}
	return
}

func Unproject(invproj mgl32.Mat4, x float32, y float32) (wx, wy float32) {
	var (
		screen = mgl32.Vec4{x, y, 1, 1}
		out    mgl32.Vec4
	)
	out = invproj.Mul4x1(screen)
	out = out.Mul(1.0 / out[3])
	wx = out[0]
	wy = out[1]
	return
}

func Project(proj mgl32.Mat4, x float32, y float32) (sx, sy float32) {
	var out mgl32.Vec4
	out = proj.Mul4x1(mgl32.Vec4{x, y, 1, 1})
	return out[0], out[1]
}
