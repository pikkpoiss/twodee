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
	gmath "github.com/Agon/googlmath"
	"github.com/go-gl/mathgl/mgl32"
)

type Matrix4 [16]float32

func (m Matrix4) Mul(a Matrix4) {
	n := getMatrix(getGMathMatrix(m).Mul(getGMathMatrix(a)))
	for i := 0; i < 16; i++ {
		m[i] = n[i]
	}
}

func getMatrix(m gmath.Matrix4) Matrix4 {
	return Matrix4{
		m.M11, m.M12, m.M13, m.M14,
		m.M21, m.M22, m.M23, m.M24,
		m.M31, m.M32, m.M33, m.M34,
		m.M41, m.M42, m.M43, m.M44,
	}
}

func getGMathMatrix(m Matrix4) gmath.Matrix4 {
	return gmath.Matrix4{
		m[0], m[1], m[2], m[3],
		m[4], m[5], m[6], m[7],
		m[8], m[9], m[10], m[11],
		m[12], m[13], m[14], m[15],
	}
}

func GetTranslationMatrix(x, y, z float32) Matrix4 {
	return getMatrix(gmath.TranslationMatrix4(x, y, z))
}

func GetRotationMatrix(x, y, z, a float32) Matrix4 {
	axis := gmath.Vector3{x, y, z}
	return getMatrix(gmath.RotationMatrix4(axis, a))
}

func GetRotTransMatrix(x, y, z, a float32) Matrix4 {
	var (
		axis  = gmath.Vector3{0, 0, 1}
		trans = gmath.TranslationMatrix4(x, y, z)
		rot   = gmath.RotationMatrix4(axis, a)
	)
	return getMatrix(trans.Mul(rot))
}

func GetRotTransScaleMatrix(x, y, z, a, s float32) Matrix4 {
	var (
		axis  = gmath.Vector3{0, 0, 1}
		trans = gmath.TranslationMatrix4(x, y, z)
		rot   = gmath.RotationMatrix4(axis, a)
		scale = getGMathMatrix(GetScaleMatrix(s, s, 0))
	)
	return getMatrix(trans.Mul(rot).Mul(scale))
}

func GetScaleMatrix(x, y, z float32) Matrix4 {
	return Matrix4{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

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
