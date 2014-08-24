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
	TexOffsetLoc   gl.UniformLocation
}

const BATCH_FRAGMENT = `#version 150
precision mediump float;

uniform sampler2D u_TextureUnit;
uniform vec2 u_TextureOffset;
in vec2 v_TextureCoordinates;
out vec4 v_FragData;

void main()
{
    vec2 texcoords = v_TextureCoordinates + u_TextureOffset;
    v_FragData = texture(u_TextureUnit, texcoords);
    //v_FragData = vec4(1.0,0.0,0.0,1.0);
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
		TexOffsetLoc:   program.GetUniformLocation("u_TextureOffset"),
	}
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (r *BatchRenderer) Bind() error {
	r.Program.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	r.TextureUnitLoc.Uniform1i(0)
	r.ProjectionLoc.UniformMatrix4f(false, (*[16]float32)(r.Renderer.projection))
	return nil
}

func (r *BatchRenderer) Draw(batch *Batch, x, y, rot float32) error {
	batch.Texture.Bind()
	batch.Buffer.Bind(gl.ARRAY_BUFFER)
	r.PositionLoc.AttribPointer(3, gl.FLOAT, false, 5*4, uintptr(0))
	r.TextureLoc.AttribPointer(2, gl.FLOAT, false, 5*4, uintptr(3*4))
	r.PositionLoc.EnableArray()
	r.TextureLoc.EnableArray()
	m := GetRotTransMatrix(x, y, 0, rot)
	r.ModelViewLoc.UniformMatrix4f(false, (*[16]float32)(m))
	r.TexOffsetLoc.Uniform2f(batch.textureOffset.X, batch.textureOffset.Y)
	gl.DrawArrays(gl.TRIANGLES, 0, batch.Count)
	batch.Buffer.Unbind(gl.ARRAY_BUFFER)
	return nil
}

func (tr *BatchRenderer) Unbind() error {
	return nil
}

func (tr *BatchRenderer) Delete() error {
	return nil
}

type Batch struct {
	Buffer        gl.Buffer
	Texture       *Texture
	Count         int
	textureOffset Point
}

type TexturedTile interface {
	ScaledBounds(ratio float32) (x, y, w, h float32)
	ScaledTextureBounds(rx float32, ry float32) (x, y, w, h float32)
}

func triangles(t TexturedTile, ratio, texw, texh float32) [30]float32 {
	var (
		x, y, w, h     = t.ScaledBounds(ratio)
		tx, ty, tw, th = t.ScaledTextureBounds(texw, texh)
	)
	return [30]float32{
		x, y, 0.0,
		tx, ty,

		x + w, y + h, 0.0,
		tx + tw, ty + th,

		x, y + h, 0.0,
		tx, ty + th,

		x, y, 0.0,
		tx, ty,

		x + w, y, 0.0,
		tx + tw, ty,

		x + w, y + h, 0.0,
		tx + tw, ty + th,
	}
}

func LoadBatch(tiles []TexturedTile, metadata TileMetadata) (b *Batch, err error) {
	var (
		step     = 30
		size     = len(tiles) * step
		vertices = make([]float32, size)
		vbo      gl.Buffer
		texture  *Texture
	)
	if texture, err = LoadTexture(metadata.Path, gl.NEAREST); err != nil {
		return
	}
	for i := 0; i < len(tiles); i++ {
		if tiles[i] == nil {
			continue
		}
		v := triangles(tiles[i], float32(metadata.PxPerUnit), float32(texture.Width), float32(texture.Height))
		copy(vertices[step*i:], v[:])
	}
	if vbo, err = CreateVBO(len(vertices)*4, vertices, gl.STATIC_DRAW); err != nil {
		return
	}
	b = &Batch{
		Buffer:        vbo,
		Texture:       texture,
		Count:         len(vertices) / 5,
		textureOffset: Pt(0, 0),
	}
	return
}

func (b *Batch) SetTextureOffsetPx(x, y int) {
	var (
		tx = float32(x) / float32(b.Texture.Width)
		ty = float32(y) / float32(b.Texture.Height)
	)
	b.textureOffset = Pt(tx, ty)
}

func (b *Batch) Delete() {
	b.Texture.Delete()
	b.Buffer.Delete()
}
