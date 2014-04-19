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

type BatchRenderer struct {
	*Renderer
	Program        gl.Program
	PositionLoc    gl.AttribLocation
	TextureLoc     gl.AttribLocation
	TextureUnitLoc gl.UniformLocation
	ModelViewLoc   gl.UniformLocation
	ProjectionLoc  gl.UniformLocation
}

const BATCH_FRAGMENT = `#version 150
precision mediump float;

uniform sampler2D u_TextureUnit;
in vec2 v_TextureCoordinates;
out vec4 v_FragData;

void main()
{
    vec2 texcoords = v_TextureCoordinates;
    v_FragData = texture(u_TextureUnit, texcoords);
}`

const BATCH_VERTEX = `#version 150

in vec4 a_Position;
in vec2 a_TextureCoordinates;

uniform mat4 m_ModelViewMatrix;
uniform mat4 m_ProjectionMatrix;

out vec2 v_TextureCoordinates;

void main()
{
    v_TextureCoordinates = a_TextureCoordinates;
    gl_Position = m_ProjectionMatrix * m_ModelViewMatrix * a_Position;
}`

func NewBatchRenderer(bounds, screen Rectangle) (tr *BatchRenderer, err error) {
	var (
		program gl.Program
		r       *Renderer
	)
	if program, err = BuildProgram(BATCH_VERTEX, BATCH_FRAGMENT); err != nil {
		return
	}
	if r, err = NewRenderer(bounds, screen); err != nil {
		return
	}
	tr = &BatchRenderer{
		Renderer:       r,
		Program:        program,
		PositionLoc:    program.GetAttribLocation("a_Position"),
		TextureLoc:     program.GetAttribLocation("a_TextureCoordinates"),
		TextureUnitLoc: program.GetUniformLocation("u_TextureUnit"),
		ModelViewLoc:   program.GetUniformLocation("m_ModelViewMatrix"),
		ProjectionLoc:  program.GetUniformLocation("m_ProjectionMatrix"),
	}
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (r *BatchRenderer) Bind() error {
	r.Program.Use()
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.ActiveTexture(gl.TEXTURE0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	r.TextureUnitLoc.Uniform1i(0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	r.ProjectionLoc.UniformMatrix4f(false, (*[16]float32)(r.Renderer.projection))
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (r *BatchRenderer) Draw(batch *Batch, x, y, rot float32) error {
	batch.Buffer.Bind(gl.ARRAY_BUFFER)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	r.PositionLoc.AttribPointer(3, gl.FLOAT, false, 5*4, uintptr(0))
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	r.TextureLoc.AttribPointer(2, gl.FLOAT, false, 5*4, uintptr(3*4))
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	r.PositionLoc.EnableArray()
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	r.TextureLoc.EnableArray()
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	m := GetRotTransMatrix(x, y, 0, rot)
	r.ModelViewLoc.UniformMatrix4f(false, (*[16]float32)(m))
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	batch.Buffer.Unbind(gl.ARRAY_BUFFER)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (tr *BatchRenderer) Unbind() error {
	return nil
}

type Batch struct {
	Buffer  gl.Buffer
	Texture *Texture
}

func LoadBatch(vertices []float32, texture *Texture) (b *Batch, err error) {
	var (
		vbo gl.Buffer
	)
	if vbo, err = CreateVBO(len(vertices)*4, vertices, gl.STATIC_DRAW); err != nil {
		return
	}
	b = &Batch{
		Buffer:  vbo,
		Texture: texture,
	}
	return
}

func (b *Batch) Delete() {
	b.Texture.Delete()
	b.Buffer.Delete()
}
