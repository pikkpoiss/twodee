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

package twodee

import (
	"github.com/jteeuwen/glfw"
	"github.com/banthar/gl"
)

type Window struct {
	Width int
	Height int
}

func (w *Window) Opened() bool {
	return glfw.WindowParam(glfw.Opened) == 1
}

type System struct {
}

func Init() (sys *System, err error) {
	if err = glfw.Init(); err != nil {
		return
	}
	sys = &System{}
	return
}

func (s *System) Key(key int) int {
	return glfw.Key(key)
}

func (s *System) Terminate() {
	glfw.Terminate()
}

func (s *System) Open(win *Window) (err error) {
	err = glfw.OpenWindow(
		win.Width,
		win.Height,
		0, 0, 0, 0 ,0, 0,
		glfw.Windowed)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(win.Width), float64(win.Height), 0, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)
	return
}

func (s *System) Paint() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Begin(gl.LINES)
	gl.Vertex2f(0, 0)
	gl.Vertex2f(640, 480)
	gl.End()
	glfw.SwapBuffers()
}
