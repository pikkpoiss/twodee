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
	"fmt"
	"log"
	"runtime"
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
	player  *Creature
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

func (g *Game) addGravity() {
	for _, e := range g.Dynamic {
		v := e.Velocity()
		e.SetVelocity(twodee.Pt(v.X, v.Y-0.001))
	}
}

func (g *Game) checkCollisions(subject twodee.Spatial) twodee.Spatial {
	b := subject.Bounds()
	for _, e := range g.Static {
		if subject == e {
			continue
		}
		if b.Overlaps(e.Bounds()) {
			return e
		}
	}
	for _, e := range g.Dynamic {
		if subject == e {
			continue
		}
		if b.Overlaps(e.Bounds()) {
			return e
		}
	}
	return nil
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

func (g *Game) checkKeys() {
	switch {
	case g.System.Key(twodee.KeyLeft) == 1 && g.System.Key(twodee.KeyRight) == 0:
		v := g.player.Velocity()
		v.X = -0.05
		g.player.SetVelocity(v)
	case g.System.Key(twodee.KeyLeft) == 0 && g.System.Key(twodee.KeyRight) == 1:
		v := g.player.Velocity()
		v.X = 0.05
		g.player.SetVelocity(v)
	}
}

func (g *Game) handleKeys() {
	g.System.SetKeyCallback(func(key int, state int) {
		switch {
		case state == 0:
			return
		case key == twodee.KeyUp:
			v := g.player.Velocity()
			v.Y = 0.1
			g.player.SetVelocity(v)
			g.ctarget.Y += 1
		case key == twodee.KeyDown:
			g.ctarget.Y -= 1
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
	g.checkKeys()
	focus := g.Camera.Focus()
	g.Camera.Pan((g.ctarget.X-focus.X)/20, (g.ctarget.Y-focus.Y)/20)
	g.addGravity()
	for _, e := range g.Dynamic {
		e.Update()
	}
	for _, e := range g.Dynamic {
		if t := g.checkCollisions(e); t != nil {
			v := e.Velocity()
			if v.Y < 0 {
				e.MoveTo(twodee.Pt(e.Bounds().Min.X, t.Bounds().Max.Y))
				e.SetVelocity(twodee.Pt(v.X, 0))
			}
		}
	}
}

func (g *Game) Create(tileset string, index int, x, y, w, h float64) {
	switch tileset {
	case "tilegame":
		var sprite = g.System.NewSprite(tileset, x, y, w, h, index)
		sprite.SetFrame(index)
		g.Static = append(g.Static, sprite)
	case "character-textures":
		var sprite = g.System.NewSprite(tileset, x, y, w, h, index)
		var creature = NewCreature(sprite)
		creature.SetFrame(index)
		g.Dynamic = append(g.Dynamic, creature)
		g.player = creature
	default:
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

const (
	STATE_STANDING = 1 << iota
	STATE_WALKING  = 1 << iota
)

type Creature struct {
	*twodee.Sprite
	Animations map[int]*twodee.Animation
	Animation  *twodee.Animation
	state      int
}

func NewCreature(sprite *twodee.Sprite) (c *Creature) {
	c = &Creature{
		Sprite: sprite,
		Animations: map[int]*twodee.Animation{
			STATE_STANDING: twodee.Anim([]int{0, 1}, 16),
			STATE_WALKING:  twodee.Anim([]int{0, 2}, 16),
		},
	}
	c.SetState(STATE_STANDING)
	return
}

func (c *Creature) SetState(state int) {
	c.state = state
	if a, ok := c.Animations[state]; ok {
		c.Animation = a
	} else {
		c.Animation = c.Animations[STATE_STANDING]
	}
}

func (c *Creature) Update() {
	if c.Sprite.VelocityX > 0 {
		c.Sprite.VelocityX -= 0.001
	} else if c.Sprite.VelocityX < 0 {
		c.Sprite.VelocityX += 0.001
	}
	c.Sprite.SetFrame(c.Animation.Next())
	c.Sprite.Update()
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
