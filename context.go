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
	"github.com/go-gl/gl/v3.3-core/gl"
	glfw "github.com/go-gl/glfw/v3.1/glfw"
)

type Context struct {
	Window        *glfw.Window
	Events        *EventHandler
	OpenGLVersion string
	ShaderVersion string
	VAO           uint32
	cursor        bool
	fullscreen    bool
	w             int
	h             int
	name          string
	initialized   bool
}

func NewContext() (context *Context, err error) {
	if err = glfw.Init(); err != nil {
		return
	}
	if err = initSound(); err != nil {
		return
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	// TODO: Understand what this implies, see
	// https://www.opengl.org/registry/specs/ARB/robustness.txt
	//glfw.WindowHint(glfw.ClientAPI, glfw.LoseContextOnReset)
	//glfw.WindowHint(glfw.ContextReleaseBehavior, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.StencilBits, 8)
	glfw.WindowHint(glfw.DepthBits, 24)
	context = &Context{
		cursor:     true,
		fullscreen: false,
	}
	return
}

func (c *Context) SetResizable(val bool) {
	if val {
		glfw.WindowHint(glfw.Resizable, 1)
	} else {
		glfw.WindowHint(glfw.Resizable, 0)
	}
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
		monitor = glfw.GetPrimaryMonitor()
	}
	if c.Window, err = glfw.CreateWindow(c.w, c.h, c.name, monitor, nil); err != nil {
		return
	}
	if c.cursor == false {
		c.Window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}
	c.Window.MakeContextCurrent()
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
	c.ShaderVersion = gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Disable(gl.CULL_FACE)
	glfw.SwapInterval(0)
	if c.VAO, err = CreateVAO(); err != nil {
		return
	}
	gl.BindVertexArray(c.VAO)
	c.Events = NewEventHandler(c.Window)
	return
}

func (c *Context) Delete() {
	cleanupSound()
	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &c.VAO)
	if c.Window != nil {
		c.Window.Destroy()
	}
	glfw.Terminate()
}
