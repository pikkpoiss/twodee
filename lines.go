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

func GetNormals(points []mgl32.Vec2, closed bool) (out []Normal) {
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
