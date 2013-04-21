// Copyright 2012 Arne Roomann-Kurrik
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
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func PrintError(err error) {
	fmt.Printf("[error]: %v\n", err)
}

func main() {
	var (
		system     *twodee.System
		window     *twodee.Window
		font       *twodee.Font
		err        error
		run        bool = true
		cpuprofile *string
		memprofile *string
		webprofile *bool
	)
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
	webprofile = flag.Bool("webprofile", false, "profile with web service")
	flag.Parse()

	if *memprofile != "" || *webprofile == true {
		runtime.MemProfileRate = 1
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *webprofile == true {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	if system, err = twodee.Init(); err != nil {
		PrintError(err)
		os.Exit(1)
	}
	defer system.Terminate()
	
	camera := twodee.NewCamera(0, 0, 10, 10)
	system.SetSizeCallback(func(w, h int) {
		camera.MatchRatio(w, h)
		camera.Top(0)
	})
	window = &twodee.Window{Width: 640, Height: 480, Scale:2}
	if err = system.Open(window); err != nil {
		PrintError(err)
		os.Exit(1)
	}

	textures := map[string]string{
		"bricks": "examples/basic/texture.png",
	}
	for name, path := range textures {
		if err = system.LoadTexture(name, path, twodee.IntNearest, 8); err != nil {
			PrintError(err)
			os.Exit(1)
		}
	}
	if font, err = twodee.LoadFont("examples/basic/slkscr.ttf", 24); err != nil {
		PrintError(err)
		os.Exit(1)
	}
	scene := &twodee.Scene{Camera: camera, Font: font}
	parent := system.NewSprite("bricks", 0, 0, 1, 1, 4)
	parent.AddChild(system.NewSprite("bricks", 1, 0.5, 1, 1, 4))
	scene.AddChild(parent)
	parent.SetFrame(1)
	exit := make(chan bool, 1)
	ticker := time.Tick(time.Second / 60.0)
	system.SetKeyCallback(func(key int, state int) {
		switch {
		case key == twodee.KeyEsc:
			exit <- true
		default:
			fmt.Printf("Key: %v, State: %v\n", key, state)
		}
	})
	system.SetCloseCallback(func() int {
		exit <- true
		return 0
	})
	system.SetScrollCallback(func(pos int) {
		fmt.Printf("Scroll: %v\n", pos)
		camera.Zoom(float64(pos) / 50.0)
		camera.Top(0)
	})
	v := twodee.Pt(0.08, 0.1)
	for run {
		worked := false
		for worked == false {
			select {
			case <-exit:
				run = false
				worked = true
			case <-ticker:
				b := parent.GlobalBounds()
				if b.Max.X > 10 || b.Min.X < 0 {
					v = twodee.Pt(-v.X, v.Y)
				} else if b.Max.Y > 10 || b.Min.Y < 0 {
					v = twodee.Pt(v.X, -v.Y)
				}
				parent.Move(v)
				worked = true
			default:
				time.Sleep(1 * time.Microsecond)
			}
		}
		system.Paint(scene)
	}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}
