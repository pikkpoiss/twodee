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

type TileMetadata struct {
	Path       string
	PxPerUnit  int
	TileWidth  int
	TileHeight int
}

type TileRenderer struct {
	*Renderer
	Program        gl.Program
	Texture        gl.Texture
	PositionLoc    gl.AttribLocation
	TextureLoc     gl.AttribLocation
	TextureUnitLoc gl.UniformLocation
	FrameLoc       gl.UniformLocation
	FramesLoc      gl.UniformLocation
	ModelViewLoc   gl.UniformLocation
	ProjectionLoc  gl.UniformLocation
	VBO            gl.Buffer
	xframes        int
	yframes        int
}

const TILE_FRAGMENT = `#version 150
precision mediump float;

uniform sampler2D u_TextureUnit;
uniform int u_Frame;
uniform ivec2 u_Frames;
in vec2 v_TextureCoordinates;
out vec4 v_FragData;

void main()
{
    vec2 scale = vec2(1.0 / float(u_Frames.x), 1.0 / float(u_Frames.y));
    vec2 texcoords = scale * v_TextureCoordinates;
    texcoords += scale * vec2(u_Frame % u_Frames.x, u_Frames.y - (u_Frame / u_Frames.x) - 1);
    v_FragData = texture(u_TextureUnit, texcoords);
}`

const TILE_VERTEX = `#version 150

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

func NewTileRenderer(bounds, screen Rectangle, metadata TileMetadata) (tr *TileRenderer, err error) {
	var (
		rect       []float32
		program    gl.Program
		texture    *Texture
		vbo        gl.Buffer
		halfWidth  = float32(metadata.TileWidth/metadata.PxPerUnit) / 2.0
		halfHeight = float32(metadata.TileHeight/metadata.PxPerUnit) / 2.0
		r          *Renderer
	)
	if program, err = BuildProgram(TILE_VERTEX, TILE_FRAGMENT); err != nil {
		return
	}
	if texture, err = LoadTexture(metadata.Path, gl.NEAREST); err != nil {
		return
	}
	rect = []float32{
		-halfWidth, -halfHeight, 0.0, 0.0, 0.0,
		-halfWidth, halfHeight, 0.0, 0.0, 1.0,
		halfWidth, -halfHeight, 0.0, 1.0, 0.0,
		halfWidth, halfHeight, 0.0, 1.0, 1.0,
	}
	if vbo, err = CreateVBO(len(rect)*4, rect, gl.STATIC_DRAW); err != nil {
		return
	}
	if r, err = NewRenderer(bounds, screen); err != nil {
		return
	}
	tr = &TileRenderer{
		Renderer:       r,
		VBO:            vbo,
		Program:        program,
		Texture:        texture.Texture,
		PositionLoc:    program.GetAttribLocation("a_Position"),
		TextureLoc:     program.GetAttribLocation("a_TextureCoordinates"),
		TextureUnitLoc: program.GetUniformLocation("u_TextureUnit"),
		FrameLoc:       program.GetUniformLocation("u_Frame"),
		FramesLoc:      program.GetUniformLocation("u_Frames"),
		ModelViewLoc:   program.GetUniformLocation("m_ModelViewMatrix"),
		ProjectionLoc:  program.GetUniformLocation("m_ProjectionMatrix"),
		xframes:        texture.Width / metadata.TileWidth,
		yframes:        texture.Height / metadata.TileHeight,
	}
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (tr *TileRenderer) Bind() error {
	tr.Program.Use()
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.ActiveTexture(gl.TEXTURE0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.Texture.Bind(gl.TEXTURE_2D)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.TextureUnitLoc.Uniform1i(0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.VBO.Bind(gl.ARRAY_BUFFER)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.PositionLoc.AttribPointer(3, gl.FLOAT, false, 5*4, uintptr(0))
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.TextureLoc.AttribPointer(2, gl.FLOAT, false, 5*4, uintptr(3*4))
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.PositionLoc.EnableArray()
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.TextureLoc.EnableArray()
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.ProjectionLoc.UniformMatrix4f(false, (*[16]float32)(tr.Renderer.projection))
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (tr *TileRenderer) Draw(frame int, x, y, r float32, flipx, flipy bool) error {
	tr.FrameLoc.Uniform1i(frame)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.FramesLoc.Uniform2i(tr.xframes, tr.yframes)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	m := GetRotTransMatrix(x, y, 0, r)
	if flipx && flipy {
		m.Mul(GetScaleMatrix(-1, -1, 1))
	} else if flipx && !flipy {
		m.Mul(GetScaleMatrix(-1, 1, 1))
	} else if !flipx && flipy {
		m.Mul(GetScaleMatrix(1, -1, 1))
	}
	tr.ModelViewLoc.UniformMatrix4f(false, (*[16]float32)(m))
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (tr *TileRenderer) Unbind() error {
	tr.VBO.Unbind(gl.ARRAY_BUFFER)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (tr *TileRenderer) Delete() error {
	tr.Texture.Delete()
	tr.VBO.Delete()
	return nil
}
