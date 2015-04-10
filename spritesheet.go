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

package twodee

import (
	"encoding/json"
	"github.com/go-gl/mathgl/mgl32"
)

type SpritesheetFrame struct {
	Points []TexturedPoint
	Frame  FrameConfig
}

type SpritesheetFrameConfig struct {
	sourceX          float32
	sourceY          float32
	sourceW          float32
	sourceH          float32
	originalW        float32
	originalH        float32
	textureX         float32
	textureY         float32
	textureW         float32
	textureH         float32
	textureOriginalW float32
	textureOriginalH float32
	pivotX           float32
	pivotY           float32
}

func (c SpritesheetFrameConfig) ToSpritesheetFrame() *SpritesheetFrame {
	var (
		x    = c.sourceX / c.originalW
		y    = c.sourceY / c.originalH
		w    = c.sourceW / c.originalW
		h    = c.sourceH / c.originalH
		texX = c.textureX / c.textureOriginalW
		texY = c.textureY / c.textureOriginalH
		texW = c.textureW / c.textureOriginalW
		texH = c.textureH / c.textureOriginalH
	)
	var (
		texMove   = mgl32.Translate3D(texX, -texH-texY, 0.0)
		texScale  = mgl32.Scale3D(texW, texH, 1.0)
		texRotate = mgl32.HomogRotate3DZ(mgl32.DegToRad(0))
		texAdj    = texMove.Mul4(texScale).Mul4(texRotate).Transpose()
	)
	var (
		ptMove   = mgl32.Translate3D(x, y, 0.0)
		ptScale  = mgl32.Scale3D(w, h, 1.0)
		ptRotate = mgl32.HomogRotate3DZ(mgl32.DegToRad(0))
		ptAdj    = ptMove.Mul4(ptScale).Mul4(ptRotate).Transpose()
	)
	return &SpritesheetFrame{
		Points: []TexturedPoint{
			TexturedPoint{x, y, 0, texX, texY - texH},
			TexturedPoint{x + w, y, 0, texX + texW, texY - texH},
			TexturedPoint{x, y + h, 0, texX, texY},
			TexturedPoint{x, y + h, 0, texX, texY},
			TexturedPoint{x + w, y, 0, texX + texW, texY - texH},
			TexturedPoint{x + w, y + h, 0, texX + texW, texY},
		},
		Frame: FrameConfig{
			PointAdjustment:   ptAdj,
			TextureAdjustment: texAdj,
		},
	}
}

type Spritesheet struct {
	frames      map[string]*SpritesheetFrame
	TexturePath string
}

func NewSpritesheet(path string) *Spritesheet {
	return &Spritesheet{
		frames:      map[string]*SpritesheetFrame{},
		TexturePath: path,
	}
}

func (s *Spritesheet) GetFrame(name string) *SpritesheetFrame {
	var (
		present bool
	)
	if _, present = s.frames[name]; !present {
		return nil
	}
	return s.frames[name]
}

func (s *Spritesheet) AddFrame(name string, config SpritesheetFrameConfig) {
	s.frames[name] = config.ToSpritesheetFrame()
}

type texturePackerFloatCoords struct {
	X float32 `json:x,omitempty`
	Y float32 `json:y,omitempty`
}

type texturePackerIntCoords struct {
	X int `json:x,omitempty`
	Y int `json:y,omitempty`
	W int `json:w,omitempty`
	H int `json:h,omitempty`
}

type texturePackerFrame struct {
	Filename         string                   `json:filename`
	Frame            texturePackerIntCoords   `json:frame`
	Rotated          bool                     `json:rotated`
	Trimmed          bool                     `json:trimmed`
	SpriteSourceSize texturePackerIntCoords   `json:spriteSourceSize`
	SourceSize       texturePackerIntCoords   `json:sourceSize`
	Pivot            texturePackerFloatCoords `json:pivot`
}

type texturePackerMeta struct {
	Image  string                 `json:image`
	Format string                 `json:format`
	Size   texturePackerIntCoords `json:size`
	Scale  string                 `json:scale`
}

type texturePackerJSONArray struct {
	Frames []texturePackerFrame `json:frames`
	Meta   texturePackerMeta    `json:meta`
}

func ParseTexturePackerJSONArrayString(contents string) (s *Spritesheet, err error) {
	var (
		parsed texturePackerJSONArray
	)
	if err = json.Unmarshal([]byte(contents), &parsed); err != nil {
		return
	}
	s = NewSpritesheet(parsed.Meta.Image)
	for _, frame := range parsed.Frames {
		s.AddFrame(frame.Filename, SpritesheetFrameConfig{
			sourceX:          float32(frame.SpriteSourceSize.X),
			sourceY:          float32(frame.SpriteSourceSize.Y),
			sourceW:          float32(frame.SpriteSourceSize.W),
			sourceH:          float32(frame.SpriteSourceSize.H),
			originalW:        float32(frame.SourceSize.W),
			originalH:        float32(frame.SourceSize.H),
			textureX:         float32(frame.Frame.X),
			textureY:         float32(frame.Frame.Y),
			textureW:         float32(frame.Frame.W),
			textureH:         float32(frame.Frame.H),
			textureOriginalW: float32(parsed.Meta.Size.W),
			textureOriginalH: float32(parsed.Meta.Size.H),
			pivotX:           frame.Pivot.X,
			pivotY:           frame.Pivot.Y,
		})
	}
	return
}
