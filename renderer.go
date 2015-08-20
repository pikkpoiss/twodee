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
	"strings"
)

type Renderer struct {
	Camera *Camera
}

func NewRenderer(camera *Camera) (r *Renderer) {
	return &Renderer{
		Camera: camera,
	}
}

func CreateVAO() (array uint32, err error) {
	gl.GenVertexArrays(1, &array)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR gl.GenVertexArray %X", e)
		return
	}
	gl.BindVertexArray(array)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR array.Bind %X", e)
		return
	}
	return
}

func CreateVBO(size int, data interface{}, usage uint32) (buffer uint32, err error) {
	gl.GenBuffers(1, &buffer)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR gl.GenBuffer %X", e)
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR buffer.Bind %X", e)
		return
	}
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(data), usage)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR gl.BufferData %X", e)
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR buffer.Unbind %X", e)
		return
	}
	return
}

func CompileShader(stype uint32, source string) (shader uint32, err error) {
	csource := gl.Str(source)
	shader = gl.CreateShader(stype)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)
	var status int32
	if gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status); status == gl.FALSE {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(shader, length, nil, gl.Str(log))
		err = fmt.Errorf("ERROR shader compile:\n%s", log)
	}
	return
}

func LinkProgram(vertex uint32, fragment uint32) (program uint32, err error) {
	program = gl.CreateProgram()
	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.BindFragDataLocation(program, 0, gl.Str("v_FragData\x00"))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR program.BindFragDataLocation %X", e)
		return
	}
	gl.LinkProgram(program)
	var status int32
	if gl.GetProgramiv(program, gl.LINK_STATUS, &status); status == gl.FALSE {
		var length int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(program, length, nil, gl.Str(log))
		err = fmt.Errorf("ERROR program link:\n%s", log)
	}
	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)
	return
}

func BuildProgram(vsrc string, fsrc string) (program uint32, err error) {
	var (
		vertex   uint32
		fragment uint32
	)
	if vertex, err = CompileShader(gl.VERTEX_SHADER, vsrc); err != nil {
		return
	}
	if fragment, err = CompileShader(gl.FRAGMENT_SHADER, fsrc); err != nil {
		return
	}
	return LinkProgram(vertex, fragment)
}
