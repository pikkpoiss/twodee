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

type TexturedPoint struct {
	X        float32
	Y        float32
	Z        float32
	TextureX float32
	TextureY float32
}

type ModelViewConfig struct {
	X         float32
	Y         float32
	Z         float32
	RotationX float32
	RotationY float32
	RotationZ float32
	ScaleX    float32
	ScaleY    float32
	ScaleZ    float32
}

type FrameConfig struct {
	PointAdjustment   mgl32.Mat4
	TextureAdjustment mgl32.Mat4
}

type SpriteConfig struct {
	View  ModelViewConfig
	Frame FrameConfig
	Color mgl32.Vec4
}
