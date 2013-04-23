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

type Map struct {
	Element
	Texture *Texture
	Blocks  []*Sprite
}

func (m *Map) Draw() {
	for _, s := range m.Blocks {
		s.Draw()
	}
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
}

type TiledTileset struct {
	Firstgid    int
	Image       string
	Imageheight int
	Imagewidth  int
	Margin      int
	Name        string
	Properties  map[string]interface{}
	Spacing     int
	Tileheight  int
	Tilewidth   int
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

func LoadTiledMap(system *System, path string) (m *Map, err error) {
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
	m = &Map{}
	for _, ts := range tm.Tilesets {
		tspath := filepath.Join(filepath.Dir(path), ts.Image)
		if err = system.LoadTexture(ts.Name, tspath, IntNearest, ts.Tilewidth); err != nil {
			return
		}
	}
	var numblocks = 0
	for _, l := range tm.Layers {
		if l.Type != "tilelayer" {
			continue
		}
		for _, f := range l.Data {
			if f == 0 {
				continue
			}
			numblocks += 1
		}
	}
	m.Blocks = make([]*Sprite, numblocks)
	var bi = 0
	for j, l := range tm.Layers {
		if l.Type != "tilelayer" {
			continue
		}
		var row, col int
		var name = tm.Tilesets[j].Name
		for i, f := range l.Data {
			if f == 0 {
				continue
			}
			row = tm.Height - 1 - i/tm.Width
			col = i % tm.Width
			sprite := system.NewSprite(name, float64(col), float64(row), 1, 1, 0)
			m.Blocks[bi] = sprite
			m.Blocks[bi].SetFrame(f - 1)
			bi += 1
		}
	}
	m.SetBounds(Rect(0, 0, float64(tm.Width), float64(tm.Height)))
	return
}
