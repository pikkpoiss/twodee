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
)

type Renderer struct {
	worldBounds  Rectangle
	screenBounds Rectangle
	projection   *Matrix4
	inverse      *Matrix4
}

func NewRenderer(world, screen Rectangle) (r *Renderer, err error) {
	r = &Renderer{}
	r.SetScreenBounds(screen)
	err = r.SetWorldBounds(world)
	return
}

func (r *Renderer) SetScreenBounds(bounds Rectangle) {
	r.screenBounds = bounds
}

func (r *Renderer) SetWorldBounds(bounds Rectangle) (err error) {
	r.worldBounds = bounds
	r.projection = GetOrthoMatrix(bounds.Min.X, bounds.Max.X, bounds.Min.Y, bounds.Max.Y, 1, 0)
	r.inverse, err = GetInverseMatrix(r.projection)
	return
}

func (r *Renderer) ScreenToWorldCoords(x, y float32) (wx, wy float32) {
	// http://stackoverflow.com/questions/7692988/
	var (
		halfw = r.screenBounds.Max.X / 2.0
		halfh = r.screenBounds.Max.Y / 2.0
		xpct  = (x - halfw) / halfw
		ypct  = (halfh - y) / halfh
	)
	return Unproject(r.inverse, xpct, ypct)
}

func CreateVAO() (array gl.VertexArray, err error) {
	array = gl.GenVertexArray()
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR gl.GenVertexArray %X", e)
		return
	}
	array.Bind()
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR array.Bind %X", e)
		return
	}
	return
}

func CreateVBO(size int, data interface{}, usage gl.GLenum) (buffer gl.Buffer, err error) {
	buffer = gl.GenBuffer()
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR gl.GenBuffer %X", e)
		return
	}
	buffer.Bind(gl.ARRAY_BUFFER)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR buffer.Bind %X", e)
		return
	}
	gl.BufferData(gl.ARRAY_BUFFER, size, data, usage)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR gl.BufferData %X", e)
		return
	}
	buffer.Unbind(gl.ARRAY_BUFFER)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR buffer.Unbind %X", e)
		return
	}
	return
}

func CompileShader(stype gl.GLenum, source string) (shader gl.Shader, err error) {
	shader = gl.CreateShader(stype)
	shader.Source(source)
	shader.Compile()
	if status := shader.Get(gl.COMPILE_STATUS); status == 0 {
		err = fmt.Errorf("ERROR shader compile:\n%s", shader.GetInfoLog())
	}
	return
}

func LinkProgram(vertex gl.Shader, fragment gl.Shader) (program gl.Program, err error) {
	program = gl.CreateProgram()
	program.AttachShader(vertex)
	program.AttachShader(fragment)
	program.BindFragDataLocation(0, "v_FragData")
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR program.BindFragDataLocation %X", e)
		return
	}
	program.Link()
	if status := program.Get(gl.LINK_STATUS); status == 0 {
		err = fmt.Errorf("ERROR program link:\n%s", program.GetInfoLog())
	}
	return
}

func BuildProgram(vsrc string, fsrc string) (program gl.Program, err error) {
	var (
		vertex   gl.Shader
		fragment gl.Shader
	)
	if vertex, err = CompileShader(gl.VERTEX_SHADER, vsrc); err != nil {
		return
	}
	if fragment, err = CompileShader(gl.FRAGMENT_SHADER, fsrc); err != nil {
		return
	}
	return LinkProgram(vertex, fragment)
}
