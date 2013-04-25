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

package main

import (
	"../.." // Use "github.com/kurrik/twodee"
	"log"
	"runtime"
	"time"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

type Updater struct {
}

func NewUpdater() *Updater {
	return &Updater{}
}

func (u *Updater) Update(e twodee.SpatialChanging) {
	e.Update()
}

type Factory struct {
	system *twodee.System
	scene  *twodee.Scene
}

func NewFactory(system *twodee.System, scene *twodee.Scene) *Factory {
	return &Factory{system: system, scene: scene}
}

func (f *Factory) SetBounds(rect twodee.Rectangle) {
	f.scene.SetBounds(rect)
}

func (f *Factory) Create(tileset string, index int, x, y, w, h float64) {
	switch tileset {
	case "tilegame":
		var sprite = f.system.NewSprite(tileset, x, y, w, h, index)
		sprite.SetFrame(index)
		f.scene.Static = append(f.scene.Static, sprite)
	default:
		var sprite = f.system.NewSprite(tileset, x, y, w, h, index)
		sprite.SetFrame(index)
		sprite.VelocityX = 0.01
		f.scene.Dynamic = append(f.scene.Dynamic, sprite)
		log.Printf("Tileset: %v %v\n", tileset, index)
		log.Printf("Dim: %v %v %v %v\n", x, y, w, h)
	}
}

func main() {
	var (
		system  *twodee.System
		camera  *twodee.Camera
		window  *twodee.Window
		factory *Factory
		updater *Updater
		font    *twodee.Font
		err     error
	)
	if system, err = twodee.Init(); err != nil {
		log.Fatalf("Couldn't init system: %v\n", err)
	}
	defer system.Terminate()

	camera = twodee.NewCamera(0, 0, 20, 20)
	cameradest := twodee.Pt(0, 0)
	system.SetSizeCallback(func(w, h int) {
		camera.MatchRatio(w, h)
		camera.Bottom(0)
	})

	window = &twodee.Window{Width: 640, Height: 480, Scale: 2}
	if err = system.Open(window); err != nil {
		log.Fatalf("Couldn't open window: %v\n", err)
	}
	system.SetClearColor(38, 147, 255, 0)
	if font, err = twodee.LoadFont("examples/complex/slkscr.ttf", 24); err != nil {
		log.Fatalf("Couldn't load font: %v\n", err)
	}
	scene := &twodee.Scene{Camera: camera, Font: font}
	factory = NewFactory(system, scene)
	if err = twodee.LoadTiledMap(system, factory, "examples/complex/levels/level01.json"); err != nil {
		log.Fatalf("Couldn't load map: %v\n", err)
	}
	updater = NewUpdater()

	exit := make(chan bool, 1)
	system.SetKeyCallback(func(key int, state int) {
		switch {
		case state == 0:
			return
		case key == twodee.KeyUp:
			cameradest.Y += 10
		case key == twodee.KeyDown:
			cameradest.Y -= 10
		case key == twodee.KeyLeft:
			cameradest.X -= 10
		case key == twodee.KeyRight:
			cameradest.X += 10
		case key == twodee.KeyEsc:
			exit <- true
		default:
			log.Printf("Key: %v, State: %v\n", key, state)
		}
	})
	system.SetCloseCallback(func() int {
		exit <- true
		return 0
	})
	lastpos := 0
	system.SetScrollCallback(func(pos int) {
		log.Printf("Scroll: %v\n", pos)
		if pos > lastpos {
			camera.Zoom(1.0 / 32.0)
		} else if pos < lastpos {
			camera.Zoom(-1.0 / 32.0)
		}
		lastpos = pos
	})
	system.SetMouseMoveCallback(func(x int, y int) {
		gx, gy := camera.ResolveScreenCoords(x, y, window.Width, window.Height)
		cameradest = twodee.Pt(gx, gy)
	})
	go func() {
		ticker := time.Tick(time.Second / 120.0)
		for true {
			<-ticker
			focus := camera.Focus()
			camera.Pan(
				(cameradest.X-focus.X)/20,
				(cameradest.Y-focus.Y)/20)
			scene.Update(updater)
		}
	}()
	ticker := time.NewTicker(time.Second / 60)
	run := true
	for run == true {
		<-ticker.C
		system.Paint(scene)
		select {
		case <-exit:
			ticker.Stop()
			run = false
		default:
		}
	}
}
