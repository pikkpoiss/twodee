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
	"fmt"
	"time"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

type Game struct {
	*twodee.Camera
	*twodee.Font
	System  *twodee.System
	Window  *twodee.Window
	Static  []twodee.SpatialVisible
	Dynamic []twodee.SpatialVisibleChanging
	exit    chan bool
	ctarget twodee.Point
	scroll  int
}

func NewGame(system *twodee.System, window *twodee.Window, font *twodee.Font) (game *Game, err error) {
	game = &Game{
		System:  system,
		Camera:  twodee.NewCamera(0, 0, 20, 20),
		Window:  window,
		Font:    font,
		exit:    make(chan bool, 1),
		scroll:  0,
		ctarget: twodee.Pt(0, 0),
	}
	game.handleResize()
	if err = system.Open(window); err != nil {
		err = fmt.Errorf("Couldn't open window: %v", err)
		return
	}
	game.handleKeys()
	game.handleClose()
	game.handleScroll()
	game.handleMouse()
	if font, err = twodee.LoadFont("examples/complex/slkscr.ttf", 24); err != nil {
		err = fmt.Errorf("Couldn't load font: %v", err)
		return
	}
	system.SetFont(font)
	system.SetClearColor(38, 147, 255, 0)
	if err = twodee.LoadTiledMap(system, game, "examples/complex/levels/level01.json"); err != nil {
		err = fmt.Errorf("Couldn't load map: %v", err)
		return
	}
	return
}

func (g *Game) handleResize() {
	g.System.SetSizeCallback(func(w, h int) {
		g.Camera.MatchRatio(w, h)
		g.Camera.Bottom(0)
	})
}

func (g *Game) handleClose() {
	g.System.SetCloseCallback(func() int {
		g.exit <- true
		return 0
	})
}

func (g *Game) handleMouse() {
	g.System.SetMouseMoveCallback(func(x int, y int) {
		gx, gy := g.Camera.ResolveScreenCoords(x, y, g.Window.Width, g.Window.Height)
		g.ctarget = twodee.Pt(gx, gy)
	})
}

func (g *Game) handleScroll() {
	g.System.SetScrollCallback(func(pos int) {
		log.Printf("Scroll: %v\n", pos)
		if pos > g.scroll {
			g.Camera.Zoom(1.0 / 32.0)
		} else if pos < g.scroll {
			g.Camera.Zoom(-1.0 / 32.0)
		}
		g.scroll = pos
	})
}

func (g *Game) handleKeys() {
	g.System.SetKeyCallback(func(key int, state int) {
		switch {
		case state == 0:
			return
		case key == twodee.KeyUp:
			g.ctarget.Y += 10
		case key == twodee.KeyDown:
			g.ctarget.Y -= 10
		case key == twodee.KeyLeft:
			g.ctarget.X -= 10
		case key == twodee.KeyRight:
			g.ctarget.X += 10
		case key == twodee.KeyEsc:
			g.exit <- true
		default:
			log.Printf("Key: %v, State: %v\n", key, state)
		}
	})
}

func (g *Game) SetBounds(rect twodee.Rectangle) {
	g.Camera.SetLimits(rect)
}

func (g *Game) Draw() {
	g.Camera.SetProjection()
	for _, e := range g.Static {
		e.Draw()
	}
	for _, e := range g.Dynamic {
		e.Draw()
	}
}

func (g *Game) Update() {
	focus := g.Camera.Focus()
	g.Camera.Pan((g.ctarget.X-focus.X)/20, (g.ctarget.Y-focus.Y)/20)
	for _, _ = range g.Dynamic {
		//e.Update()
	}
}

func (g *Game) Create(tileset string, index int, x, y, w, h float64) {
	switch tileset {
	case "tilegame":
		var sprite = g.System.NewSprite(tileset, x, y, w, h, index)
		sprite.SetFrame(index)
		g.Static = append(g.Static, sprite)
	default:
		var sprite = g.System.NewSprite(tileset, x, y, w, h, index)
		sprite.SetFrame(index)
		sprite.VelocityX = 0.01
		g.Dynamic = append(g.Dynamic, sprite)
		log.Printf("Tileset: %v %v\n", tileset, index)
		log.Printf("Dim: %v %v %v %v\n", x, y, w, h)
	}
}

func (g *Game) Run() {
	go func() {
		ticker := time.Tick(time.Second / 120.0)
		for true {
			<-ticker
			g.Update()
		}
	}()
	ticker := time.NewTicker(time.Second / 60)
	run := true
	for run == true {
		<-ticker.C
		g.System.Paint(g)
		select {
		case <-g.exit:
			ticker.Stop()
			run = false
		default:
		}
	}
}

func main() {
	var (
		system *twodee.System
		window *twodee.Window
		game   *Game
		font   *twodee.Font
		err    error
	)
	if system, err = twodee.Init(); err != nil {
		log.Fatalf("Couldn't init system: %v\n", err)
	}
	defer system.Terminate()
	window = &twodee.Window{Width: 640, Height: 480, Scale: 2}
	if game, err = NewGame(system, window, font); err != nil {
		log.Fatalf("Couldn't run game: %v\n", err)
	}
	game.Run()
}
