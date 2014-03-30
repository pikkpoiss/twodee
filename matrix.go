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

func GetMatrix(m *gmath.Matrix4) *[16]float32 {
	return &[16]float32{
		m.M11, m.M12, m.M13, m.M14,
		m.M21, m.M22, m.M23, m.M24,
		m.M31, m.M32, m.M33, m.M34,
		m.M41, m.M42, m.M43, m.M44,
	}
}

func GetTranslationMatrix(x, y, z float32) *[16]float32 {
	return GetMatrix(gmath.NewTranslationMatrix4(x, y, z))
}

func GetRotationMatrix(x, y, z, a float32) *[16]float32 {
	axis := gmath.Vector3{x, y, z}
	return GetMatrix(gmath.NewRotationMatrix4(axis, a))
}

func GetRotTransMatrix(x, y, z, a float32) *[16]float32 {
	var (
		axis  = gmath.Vector3{0, 0, 1}
		trans = gmath.NewTranslationMatrix4(x, y, z)
		rot   = gmath.NewRotationMatrix4(axis, a)
	)
	return GetMatrix(trans.Mul(rot))
}

func GetScaleMatrix(x, y, z float32) *[16]float32 {
	return &[16]float32{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

func GetOrthoMatrix(x1, x2, y1, y2, n, f float32) *[16]float32 {
	// http://www.songho.ca/opengl/gl_projectionmatrix.html
	return &[16]float32{
		2.0 / (x2 - x1), 0, 0, 0,
		0, 2.0 / (y2 - y1), 0, 0,
		0, 0, -2.0 / (f - n), 0,
		-(x2 + x1) / (x2 - x1), -(y2 + y1) / (y2 - y1), -(f + n) / (f - n), 1,
	}
}


