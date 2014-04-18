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
	if program, err = BuildProgram(TILE_VERTEX, TILE_FRAGMENT); err != nil {
		return
	}
	if r, err = NewRenderer(bounds, screen); err != nil {
		return
	}
	tr = &BatchRenderer{
		Renderer:       r,
		Program:        program,
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
