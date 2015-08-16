// Copyright 2015 Arne Roomann-Kurrik
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

// Ported from
//   https://github.com/mattdesl/polyline-normals
//   https://github.com/mattdesl/polyline-miter-util
//   https://github.com/stackgl/gl-vec2

package twodee

import (
	"github.com/go-gl/mathgl/mgl32"
)

func computeMiter(lineA, lineB mgl32.Vec2, halfThick float32) (miter mgl32.Vec2, length float32) {
	var (
		tangent mgl32.Vec2
		tmp     mgl32.Vec2
	)
	tangent = lineA.Add(lineB).Normalize()
	miter = mgl32.Vec2{-tangent[1], tangent[0]}
	tmp = mgl32.Vec2{-lineA[1], lineA[0]}
	length = halfThick / miter.Dot(tmp)
	return
}

func normal(dir mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{-dir[1], dir[0]}
}

func direction(a, b mgl32.Vec2) mgl32.Vec2 {
	return a.Sub(b).Normalize()
}

func getNormals(points []mgl32.Vec2, closed bool) (out []Normal) {
	var (
		curNormal    mgl32.Vec2
		total        int
		i            int
		last         mgl32.Vec2
		cur          mgl32.Vec2
		next         mgl32.Vec2
		lineA        mgl32.Vec2
		lineB        mgl32.Vec2
		miterLen     float32
		miter        mgl32.Vec2
		hasNext      bool = false
		hasCurNormal bool = false
	)
	out = make([]Normal, 0)
	if closed {
		points = append(points, points[0])
	}
	total = len(points)
	for i = 1; i < total; i++ {
		last = points[i-1]
		cur = points[i]
		if i < (len(points) - 1) {
			next = points[i+1]
			hasNext = true
		} else {
			hasNext = false
		}
		lineA = direction(cur, last)
		if !hasCurNormal {
			curNormal = normal(lineA)
			hasCurNormal = true
		}
		if i == 1 {
			out = addNext(out, curNormal, 1)
		}
		if !hasNext {
			curNormal = normal(lineA)
			out = addNext(out, curNormal, 1)
		} else {
			lineB = direction(next, cur)
			miter, miterLen = computeMiter(lineA, lineB, 1)
			out = addNext(out, miter, miterLen)
		}
	}
	if len(points) > 2 && closed {
		var (
			last2 = points[total-2]
			cur2  = points[0]
			next2 = points[1]
		)
		lineA = direction(cur2, last2)
		lineB = direction(next2, cur2)
		curNormal = normal(lineA)
		miter, miterLen = computeMiter(lineA, lineB, 1)
		out[0] = Normal{Vector: miter, Length: miterLen}
		out = out[:len(out)-1]
	}
	return
}

type Normal struct {
	Vector mgl32.Vec2
	Length float32
}

func addNext(list []Normal, normal mgl32.Vec2, length float32) []Normal {
	return append(list, Normal{Vector: normal, Length: length})
}

func duplicateNormals(list []Normal) (out []Normal) {
	out = make([]Normal, len(list)*2)
	for i := 0; i < len(list); i++ {
		out[2*i] = list[i]
		out[2*i].Length *= -1
		out[2*i+1] = list[i]
	}
	return
}

func duplicateVec2(list []mgl32.Vec2) (out []mgl32.Vec2) {
	out = make([]mgl32.Vec2, len(list)*2)
	for i := 0; i < len(list); i++ {
		out[2*i] = list[i]
		out[2*i+1] = list[i]
	}
	return
}

type LineGeometry struct {
	Points   []TexturedPoint
	Vertices []TexturedPoint
	Indices  []uint32
}

func NewLineGeometry(path []mgl32.Vec2, closed bool) (out *LineGeometry) {
	var (
		normals  []Normal
		count    int
		indices  []uint32
		vertices []TexturedPoint
		i        int
	)
	normals = getNormals(path, closed)
	if closed {
		normals = append(normals, normals[0])
		path = append(path, path[0])
	}
	count = len(path) - 1
	indices = make([]uint32, count*6)
	for i = 0; i < count; i++ {
		indices[i*6+0] = uint32(2*i + 0)
		indices[i*6+1] = uint32(2*i + 1)
		indices[i*6+2] = uint32(2*i + 2)
		indices[i*6+3] = uint32(2*i + 2)
		indices[i*6+4] = uint32(2*i + 1)
		indices[i*6+5] = uint32(2*i + 3)
	}
	normals = duplicateNormals(normals)
	path = duplicateVec2(path)
	vertices = make([]TexturedPoint, len(normals))
	for i = 0; i < len(normals); i++ {
		vertices[i] = TexturedPoint{
			X:        path[i][0],
			Y:        path[i][1],
			Z:        normals[i].Length,
			TextureX: normals[i].Vector[0],
			TextureY: normals[i].Vector[1],
		}
	}
	out = &LineGeometry{
		Indices:  indices,
		Vertices: vertices,
	}
	return
}
