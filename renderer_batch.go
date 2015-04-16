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
	"github.com/go-gl/mathgl/mgl32"
)

type InterpolationType uint32

const (
	LinearInterpolation  = InterpolationType(gl.LINEAR)
	NearestInterpolation = InterpolationType(gl.NEAREST)
)

type TileMetadata struct {
	Path          string
	PxPerUnit     int
	TileWidth     int
	TileHeight    int
	FramesWide    int
	FramesHigh    int
	Interpolation InterpolationType
}

type BatchRenderer struct {
	*Renderer
	Program        uint32
	PositionLoc    uint32
	TextureLoc     uint32
	TextureUnitLoc int32
	ModelViewLoc   int32
	ProjectionLoc  int32
	TexOffsetLoc   int32
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
}` + "\x00"

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
}` + "\x00"

func NewBatchRenderer(bounds, screen Rectangle) (tr *BatchRenderer, err error) {
	var (
		program uint32
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
		PositionLoc:    uint32(gl.GetAttribLocation(program, gl.Str("a_Position\x00"))),
		TextureLoc:     uint32(gl.GetAttribLocation(program, gl.Str("a_TextureCoordinates\x00"))),
		TextureUnitLoc: gl.GetUniformLocation(program, gl.Str("u_TextureUnit\x00")),
		ModelViewLoc:   gl.GetUniformLocation(program, gl.Str("m_ModelViewMatrix\x00")),
		ProjectionLoc:  gl.GetUniformLocation(program, gl.Str("m_ProjectionMatrix\x00")),
		TexOffsetLoc:   gl.GetUniformLocation(program, gl.Str("u_TextureOffset\x00")),
	}
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (r *BatchRenderer) Bind() error {
	gl.UseProgram(r.Program)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(r.TextureUnitLoc, 0)
	gl.UniformMatrix4fv(r.ProjectionLoc, 1, false, &r.Renderer.projection[0])
	return nil
}

func (r *BatchRenderer) Draw(batch *Batch, x, y, rot float32) error {
	batch.Texture.Bind()
	gl.BindBuffer(gl.ARRAY_BUFFER, batch.Buffer)
	gl.EnableVertexAttribArray(r.PositionLoc)
	gl.VertexAttribPointer(r.PositionLoc, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(r.TextureLoc)
	gl.VertexAttribPointer(r.TextureLoc, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	m := mgl32.Translate3D(x, y, 0.0).Mul4(mgl32.HomogRotate3DZ(rot))
	gl.UniformMatrix4fv(r.ModelViewLoc, 1, false, &m[0])
	gl.Uniform2f(r.TexOffsetLoc, batch.textureOffset.X, batch.textureOffset.Y)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(batch.Count))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return nil
}

func (tr *BatchRenderer) Unbind() error {
	return nil
}

func (tr *BatchRenderer) Delete() error {
	return nil
}

type Batch struct {
	Buffer        uint32
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
		vbo      uint32
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
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &b.Buffer)
}
