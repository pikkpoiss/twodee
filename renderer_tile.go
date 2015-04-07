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

type TileRenderer struct {
	*Renderer
	Program        uint32
	Texture        uint32
	PositionLoc    uint32
	TextureLoc     uint32
	TextureUnitLoc int32
	FrameLoc       int32
	FramesLoc      int32
	RatioLoc       int32
	ModelViewLoc   int32
	ProjectionLoc  int32
	VBO            uint32
	xframes        int
	yframes        int
	xratio         float32
	yratio         float32
}

const TILE_FRAGMENT = `#version 150
precision mediump float;

uniform sampler2D u_TextureUnit;
uniform int u_Frame;
uniform ivec2 u_Frames;
uniform vec2 u_Ratio;
in vec2 v_TextureCoordinates;
out vec4 v_FragData;

void main()
{
    vec2 scale = vec2(u_Ratio.x / float(u_Frames.x), u_Ratio.y / float(u_Frames.y));
    vec2 texcoords = scale * v_TextureCoordinates;
    texcoords += scale * vec2(u_Frame % u_Frames.x, u_Frames.y - (u_Frame / u_Frames.x) - 1);
    texcoords.y += (1 - u_Ratio.y);
    vec4 color = texture(u_TextureUnit, texcoords);
    if (color.a < 0.1) {
      discard;
    }
    v_FragData = color;
}` + "\x00"

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
}` + "\x00"

func NewTileRenderer(bounds, screen Rectangle, metadata TileMetadata) (tr *TileRenderer, err error) {
	var (
		texRatioX  float32
		texRatioY  float32
		rect       []float32
		program    uint32
		texture    *Texture
		vbo        uint32
		halfWidth  = float32(metadata.TileWidth/metadata.PxPerUnit) / 2.0
		halfHeight = float32(metadata.TileHeight/metadata.PxPerUnit) / 2.0
		r          *Renderer
	)
	if program, err = BuildProgram(TILE_VERTEX, TILE_FRAGMENT); err != nil {
		return
	}
	if metadata.Interpolation == 0 {
		metadata.Interpolation = NearestInterpolation
	}
	if texture, err = LoadTexture(metadata.Path, int(metadata.Interpolation)); err != nil {
		return
	}
	texRatioX = float32(texture.OriginalWidth) / float32(texture.Width)
	texRatioY = float32(texture.OriginalHeight) / float32(texture.Height)
	rect = []float32{
		-halfWidth, -halfHeight, 0.0, 0, 0,
		-halfWidth, halfHeight, 0.0, 0, 1,
		halfWidth, -halfHeight, 0.0, 1, 0,
		halfWidth, halfHeight, 0.0, 1, 1,
	}
	if vbo, err = CreateVBO(len(rect)*4, rect, gl.STATIC_DRAW); err != nil {
		return
	}
	if r, err = NewRenderer(bounds, screen); err != nil {
		return
	}
	if metadata.FramesWide == 0 {
		metadata.FramesWide = texture.OriginalWidth / metadata.TileWidth
	}
	if metadata.FramesHigh == 0 {
		metadata.FramesHigh = texture.OriginalHeight / metadata.TileHeight
	}
	tr = &TileRenderer{
		Renderer:       r,
		VBO:            vbo,
		Program:        program,
		Texture:        texture.Texture,
		PositionLoc:    uint32(gl.GetAttribLocation(program, gl.Str("a_Position\x00"))),
		TextureLoc:     uint32(gl.GetAttribLocation(program, gl.Str("a_TextureCoordinates\x00"))),
		TextureUnitLoc: gl.GetUniformLocation(program, gl.Str("u_TextureUnit\x00")),
		FrameLoc:       gl.GetUniformLocation(program, gl.Str("u_Frame\x00")),
		FramesLoc:      gl.GetUniformLocation(program, gl.Str("u_Frames\x00")),
		RatioLoc:       gl.GetUniformLocation(program, gl.Str("u_Ratio\x00")),
		ModelViewLoc:   gl.GetUniformLocation(program, gl.Str("m_ModelViewMatrix\x00")),
		ProjectionLoc:  gl.GetUniformLocation(program, gl.Str("m_ProjectionMatrix\x00")),
		xframes:        metadata.FramesWide,
		yframes:        metadata.FramesHigh,
		xratio:         texRatioX,
		yratio:         texRatioY,
	}
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (tr *TileRenderer) Bind() error {
	gl.UseProgram(tr.Program)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tr.Texture)
	gl.Uniform1i(tr.TextureUnitLoc, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, tr.VBO)
	gl.EnableVertexAttribArray(tr.PositionLoc)
	gl.VertexAttribPointer(tr.PositionLoc, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(tr.TextureLoc)
	gl.VertexAttribPointer(tr.TextureLoc, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.UniformMatrix4fv(tr.ProjectionLoc, 1, false, &tr.Renderer.projection[0])
	return nil
}

func (tr *TileRenderer) Draw(frame int, x, y, r float32, flipx, flipy bool) error {
	gl.Uniform1i(tr.FrameLoc, int32(frame))
	gl.Uniform2i(tr.FramesLoc, int32(tr.xframes), int32(tr.yframes))
	gl.Uniform2f(tr.RatioLoc, tr.xratio, tr.yratio)
	m := GetRotTransMatrix(x, y, 0, r)
	if flipx && flipy {
		m.Mul(GetScaleMatrix(-1, -1, 1))
	} else if flipx && !flipy {
		m.Mul(GetScaleMatrix(-1, 1, 1))
	} else if !flipx && flipy {
		m.Mul(GetScaleMatrix(1, -1, 1))
	}
	gl.UniformMatrix4fv(tr.ModelViewLoc, 1, false, &m[0])
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	return nil
}

func (tr *TileRenderer) DrawScaled(frame int, x, y, r, s float32, flipx, flipy bool) error {
	gl.Uniform1i(tr.FrameLoc, int32(frame))
	gl.Uniform2i(tr.FramesLoc, int32(tr.xframes), int32(tr.yframes))
	gl.Uniform2f(tr.RatioLoc, tr.xratio, tr.yratio)
	m := GetRotTransScaleMatrix(x, y, 0, r, s)
	if flipx && flipy {
		m.Mul(GetScaleMatrix(-1, -1, 1))
	} else if flipx && !flipy {
		m.Mul(GetScaleMatrix(-1, 1, 1))
	} else if !flipx && flipy {
		m.Mul(GetScaleMatrix(1, -1, 1))
	}
	gl.UniformMatrix4fv(tr.ModelViewLoc, 1, false, &m[0])
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	return nil
}

func (tr *TileRenderer) Unbind() error {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (tr *TileRenderer) Delete() error {
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.DeleteTextures(1, &tr.Texture)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &tr.VBO)
	return nil
}
