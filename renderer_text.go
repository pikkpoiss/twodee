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
)

type TextRenderer struct {
	*Renderer
	VBO            uint32
	Program        uint32
	Texture        uint32
	PositionLoc    uint32
	TextureLoc     uint32
	ScaleLoc       int32
	TransLoc       int32
	ProjectionLoc  int32
	TextureUnitLoc int32
	projection     *Matrix4
	Width          float32
	Height         float32
}

const TEXT_FRAGMENT = `#version 150
precision mediump float;

uniform sampler2D u_TextureUnit;
in vec2 v_TextureCoordinates;
out vec4 v_FragData;

void main()
{
    v_FragData = texture(u_TextureUnit, v_TextureCoordinates);
}` + "\x00"

const TEXT_VERTEX = `#version 150

in vec4 a_Position;
in vec2 a_TextureCoordinates;

uniform mat4 m_ProjectionMatrix;
uniform vec3 v_Trans;
uniform vec3 v_Scale;
out vec2 v_TextureCoordinates;

void main()
{
    mat4 trans;
    trans[0] = vec4(1,0,0,0);
    trans[1] = vec4(0,1,0,0);
    trans[2] = vec4(0,0,1,0);
    trans[3] = vec4(v_Trans.x,v_Trans.y,v_Trans.z,1);

    mat4 scale;
    scale[0] = vec4(v_Scale.x,0,0,0);
    scale[1] = vec4(0,v_Scale.y,0,0);
    scale[2] = vec4(0,0,v_Scale.z,0);
    scale[3] = vec4(0,0,0,1);

    v_TextureCoordinates = a_TextureCoordinates;
    gl_Position = m_ProjectionMatrix * trans * scale * a_Position;
}` + "\x00"

func NewTextRenderer(screen Rectangle) (tr *TextRenderer, err error) {
	var (
		rect    []float32
		program uint32
		vbo     uint32
		r       *Renderer
	)
	rect = []float32{
		0, 0, 0.0, 0.0, 0.0,
		0, 1, 0.0, 0.0, 1.0,
		1, 0, 0.0, 1.0, 0.0,
		1, 1, 0.0, 1.0, 1.0,
	}
	if program, err = BuildProgram(TEXT_VERTEX, TEXT_FRAGMENT); err != nil {
		return
	}
	if vbo, err = CreateVBO(len(rect)*4, rect, gl.STATIC_DRAW); err != nil {
		return
	}
	if r, err = NewRenderer(screen, screen); err != nil {
		return
	}
	tr = &TextRenderer{
		Renderer:       r,
		VBO:            vbo,
		Program:        program,
		PositionLoc:    uint32(gl.GetAttribLocation(program, gl.Str("a_Position\x00"))),
		TextureLoc:     uint32(gl.GetAttribLocation(program, gl.Str("a_TextureCoordinates\x00"))),
		TextureUnitLoc: gl.GetUniformLocation(program, gl.Str("u_TextureUnit\x00")),
		TransLoc:       gl.GetUniformLocation(program, gl.Str("v_Trans\x00")),
		ScaleLoc:       gl.GetUniformLocation(program, gl.Str("v_Scale\x00")),
		ProjectionLoc:  gl.GetUniformLocation(program, gl.Str("m_ProjectionMatrix\x00")),
		Width:          screen.Max.X - screen.Min.X,
		Height:         screen.Max.Y - screen.Min.Y,
	}
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (tr *TextRenderer) Bind() error {
	gl.UseProgram(tr.Program)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.Uniform1i(tr.TextureUnitLoc, 0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, tr.VBO)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.EnableVertexAttribArray(tr.PositionLoc)
	gl.VertexAttribPointer(tr.PositionLoc, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.EnableVertexAttribArray(tr.TextureLoc)
	gl.VertexAttribPointer(tr.TextureLoc, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.UniformMatrix4fv(tr.ProjectionLoc, 1, false, &tr.Renderer.projection[0])
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (tr *TextRenderer) Draw(tex *Texture, x, y float32) (err error) {
	gl.ActiveTexture(gl.TEXTURE0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.BindTexture(gl.TEXTURE_2D, tex.Texture)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.Uniform3f(tr.ScaleLoc, float32(tex.Width), float32(tex.Height), 1)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.Uniform3f(tr.TransLoc, x, y, 0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (tr *TextRenderer) Unbind() error {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (tr *TextRenderer) Delete() error {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &tr.VBO)
	return nil
}
