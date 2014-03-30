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
	OpenGLVersion string
	ShaderVersion string
	VAO           gl.VertexArray
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
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.ClientApi, glfw.OpenglApi)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	glfw.WindowHint(glfw.OpenglForwardCompatible, gl.TRUE)
	context = &Context{}
	return
}

func (c *Context) CreateWindow(w, h int, name string) (err error) {
	if c.Window, err = glfw.CreateWindow(w, h, name, nil, nil); err != nil {
		return
	}
	c.Window.MakeContextCurrent()
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("OpenGL MakeContextCurrent error: %X\n", e)
		return
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
	return
}

func (c *Context) Delete() {
	c.VAO.Delete()
	glfw.Terminate()
}
