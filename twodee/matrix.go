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
	"github.com/go-gl/mathgl/mgl32"
)

func GetTranslationMatrix(x, y, z float32) mgl32.Mat4 {
	return mgl32.Translate3D(x, y, z)
}

func GetRotationMatrix(x, y, z, a float32) mgl32.Mat4 {
	axis := mgl32.Vec3{x, y, z}
	return HomogRotate3D(a, axis)
}

func GetScaleMatrix(x, y, z float32) mgl32.Mat4 {
	return Scale3D(x, y, z)
}

func GetRotTransMatrix(x, y, z, a float32) mgl32.Mat4 {
	var (
		trans = GetTranslationMatrix(x, y, z)
		rot   = GetRotationMatrix(0, 0, 1, a)
	)
	return trans.Mul(rot)
}

func GetRotTransScaleMatrix(x, y, z, a, s float32) mgl32.Mat4 {
	var (
		trans = GetTranslationMatrix(x, y, z)
		rot   = GetRotationMatrix(0, 0, 1, a)
		scale = GetScaleMatrix(s, s, 1)
	)
	return trans.Mul(rot).Mul(scale)
}

func GetOrthoMatrix(x1, x2, y1, y2, n, f float32) mgl32.Mat4 {
	return mgl32.Ortho(x1, x2, y1, y2, near, far)
}

func GetInverseMatrix(m mgl32.Mat4) (out mgl32.Mat4, err error) {
	var inv = m.Inv()
	if inv == mgl.Mat4 {
		err = fmt.Errorf("Cannot invert matrix")
	}
	return
}

func Unproject(invproj mgl32.Mat4, x float32, y float32) (wx, wy float32) {
	var (
		screen = mgl32.Vec4{x, y, 1, 1}
		out    mgl32.Vec4
	)
	out = invproj.Mul4x1(screen).Mul(1.0 / out.W)
	wx = out.X()
	wy = out.Y()
	return
}

func Project(proj mgl32.Mat4, x float32, y float32) (sx, sy float32) {
	var out mgl32.Vec4
	out = proj.Mul4x1(mgl32.Vec4{x, y, 1, 1})
	return out.X, out.Y
}
