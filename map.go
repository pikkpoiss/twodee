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

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type TiledObject struct {
	Height     int
	Name       string
	Properties map[string]interface{}
	Type       string
	Width      int
	X          int
	Y          int
}

type TiledLayer struct {
	Data    []int
	Height  int
	Name    string
	Opacity float32
	Type    string
	Visible bool
	Width   int
	X       int
	Y       int
	Objects []TiledObject
}

type TiledTileset struct {
	Firstgid         int
	Image            string
	Imageheight      int
	Imagewidth       int
	Margin           int
	Name             string
	Properties       map[string]interface{}
	Spacing          int
	Tileheight       int
	Tilewidth        int
	Transparentcolor string
}

type TiledMap struct {
	Height      int
	Layers      []TiledLayer
	Orientation string
	Properties  map[string]interface{}
	Tileheight  int
	Tilesets    []TiledTileset
	Tilewidth   int
	Version     int
	Width       int
}

func LoadTiledMap(system *System, loader MapLoader, path string) (err error) {
	var (
		f       *os.File
		decoder *json.Decoder
		tm      TiledMap
	)
	if f, err = os.Open(path); err != nil {
		return
	}
	defer f.Close()
	decoder = json.NewDecoder(f)
	if err = decoder.Decode(&tm); err != nil {
		return
	}
	var gids = make([]int, len(tm.Tilesets))
	for i, ts := range tm.Tilesets {
		tspath := filepath.Join(filepath.Dir(path), ts.Image)
		if err = system.LoadTexture(ts.Name, tspath, IntNearest, ts.Tilewidth); err != nil {
			return
		}
		gids[i] = ts.Firstgid
	}
	for _, l := range tm.Layers {
		var row, col, width, height float64
		switch l.Type {
		case "objectgroup":
			for _, o := range l.Objects {
				col = float64(o.X) / float64(tm.Tilewidth)
				row = float64(o.Y) / float64(tm.Tileheight)
				width = float64(o.Width) / float64(tm.Tilewidth)
				height = float64(o.Height) / float64(tm.Tileheight)
				loader.Create(o.Name, -1, col, row, width, height)
			}
		case "tilelayer":
			var ts TiledTileset
			for i, f := range l.Data {
				if f == 0 {
					continue
				}
				var tsi = len(gids) - 1
				for gids[tsi] > f {
					tsi -= 1
				}
				f = f - gids[tsi]
				row = float64(tm.Height - 1 - i/tm.Width)
				col = float64(i % tm.Width)
				ts = tm.Tilesets[tsi]
				height = float64(ts.Tileheight) / float64(tm.Tileheight)
				width = float64(ts.Tilewidth) / float64(tm.Tilewidth)
				loader.Create(ts.Name, f, col, row, width, height)
			}
		}
	}
	loader.SetBounds(Rect(0, 0, float64(tm.Width), float64(tm.Height)))
	return
}
