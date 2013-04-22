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
	"fmt"
	"os"
)

type Map struct {
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

func LoadTiledMap(path string) (m *Map, err error) {
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
	fmt.Printf("%v\n", tm)
	return
}
