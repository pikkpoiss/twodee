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
	"testing"
)

const TEST_ARRAY_STRING = `{
  "frames": [
    {
      "filename": "frame01.png",
      "frame": {"x":2,"y":112,"w":26,"h":26},
      "rotated": false,
      "trimmed": true,
      "spriteSourceSize": {"x":3,"y":3,"w":26,"h":26},
      "sourceSize": {"w":32,"h":32},
      "pivot": {"x":0.5,"y":0.5}
    },
    {
      "filename": "frame02.png",
      "frame": {"x":2,"y":2,"w":30,"h":30},
      "rotated": false,
      "trimmed": true,
      "spriteSourceSize": {"x":1,"y":1,"w":30,"h":30},
      "sourceSize": {"w":32,"h":32},
      "pivot": {"x":0.5,"y":0.5}
    }
  ],
  "meta": {
    "app": "http://www.codeandweb.com/texturepacker",
    "version": "1.0",
    "image": "test.png",
    "format": "RGBA8888",
    "size": {"w":64,"h":271},
    "scale": "1",
    "smartupdate": "$TexturePacker:SmartUpdate:xxxx"
  }
}`

func TestParseTexturePackerJSONArrayString(t *testing.T) {
	var (
		sheet  *Spritesheet
		frame  *SpritesheetFrame
		err    error
		point0 = TexturedPoint{-0.40625, -0.40625, 0.0, 0.03125, 0.41328412}
		point1 = TexturedPoint{-0.40625, 0.40625, 0.0, 0.03125, 0.5092251}
		point2 = TexturedPoint{0.40625, -0.40625, 0.0, 0.4375, 0.41328412}
		point3 = TexturedPoint{0.40625, 0.40625, 0.0, 0.4375, 0.5092251}
	)
	if sheet, err = ParseTexturePackerJSONArrayString(TEST_ARRAY_STRING); err != nil {
		t.Fatalf("Problem parsing JSON array: %v", err)
	}
	if sheet.GetFrame("non_existent_frame.png") != nil {
		t.Fatalf("GetFrame with unknown name should return nil")
	}
	if frame = sheet.GetFrame("frame01.png"); frame == nil {
		t.Fatalf("GetFrame with known name should return non-nil")
	}
	if frame.Points[0] != point0 {
		t.Fatalf("Invalid point, got %v, expected %v", frame.Points[0], point0)
	}
	if frame.Points[1] != point1 {
		t.Fatalf("Invalid point, got %v, expected %v", frame.Points[1], point1)
	}
	if frame.Points[2] != point2 {
		t.Fatalf("Invalid point, got %v, expected %v", frame.Points[2], point2)
	}
	if frame.Points[3] != point3 {
		t.Fatalf("Invalid point, got %v, expected %v", frame.Points[3], point3)
	}
}