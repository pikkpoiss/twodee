// Copyright 2014 Arne Roomann-Kurrik
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
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Context struct {
	Window        *glfw.Window
	Events        *EventHandler
	OpenGLVersion string
	ShaderVersion string
	VAO           gl.VertexArray
	cursor        bool
	fullscreen    bool
	w             int
	h             int
	name          string
	initialized   bool
}

func glfwErrorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func NewContext() (context *Context, err error) {
	glfw.SetErrorCallback(glfwErrorCallback)
	if !glfw.Init() {
		err = fmt.Errorf("Could not init glfw")
		return
	}
	if err = initSound(); err != nil {
		return
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.ClientApi, glfw.OpenglApi)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	glfw.WindowHint(glfw.OpenglForwardCompatible, gl.TRUE)
	context = &Context{
		cursor:     true,
		fullscreen: false,
	}
	return
}

func (c *Context) SetCursor(val bool) {
	c.cursor = val
}

func (c *Context) SetFullscreen(val bool) {
	if c.fullscreen == val {
		return
	}
	c.fullscreen = val
	if c.Window != nil {
		win := c.Window
		c.Window = nil
		win.Destroy()
		c.CreateWindow(c.w, c.h, c.name)
	}
}

func (c *Context) Fullscreen() bool {
	return c.fullscreen
}

func (c *Context) CreateWindow(w, h int, name string) (err error) {
	c.w = w
	c.h = h
	c.name = name
	var monitor *glfw.Monitor
	if c.fullscreen == true {
		if monitor, err = glfw.GetPrimaryMonitor(); err != nil {
			return
		}
	}
	if c.Window, err = glfw.CreateWindow(c.w, c.h, c.name, monitor, nil); err != nil {
		return
	}
	if c.cursor == false {
		c.Window.SetInputMode(glfw.Cursor, glfw.CursorHidden)
	}
	c.Window.MakeContextCurrent()
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("OpenGL MakeContextCurrent error: %X\n", e)
	}
	gl.Init()
	if e := gl.GetError(); e != 0 {
		if e == gl.INVALID_ENUM {
			fmt.Printf("GL_INVALID_ENUM when calling glInit\n")
		} else {
			err = fmt.Errorf("OpenGL glInit error: %X\n", e)
			return
		}
	}
	c.OpenGLVersion = glfw.GetVersionString()
	c.ShaderVersion = gl.GetString(gl.SHADING_LANGUAGE_VERSION)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Disable(gl.CULL_FACE)
	glfw.SwapInterval(1)
	if c.VAO, err = CreateVAO(); err != nil {
		return
	}
	c.VAO.Bind()
	c.Events = NewEventHandler(c.Window)
	return
}

func (c *Context) Delete() {
	cleanupSound()
	c.VAO.Delete()
	if c.Window != nil {
		c.Window.Destroy()
	}
	glfw.Terminate()
}
