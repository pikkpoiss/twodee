// Copyright 2015 Arne Roomann-Kurrik
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
	"image/color"
	"unsafe"
)

type LinesRenderer struct {
	*Renderer
	program       uint32
	buffer        uint32
	indexBuffer   uint32
	bufferBytes   int
	positionLoc   uint32
	normalLoc     uint32
	miterLoc      uint32
	modelviewLoc  int32
	thicknessLoc  int32
	colorLoc      int32
	innerLoc      int32
	projectionLoc int32
	offPosition   unsafe.Pointer
	offNormal     unsafe.Pointer
	offMiter      unsafe.Pointer
	stride        int32
}

const LINES_FRAGMENT = `#version 150
precision mediump float;

in float f_Edge;
uniform float f_Inner;
uniform vec4 v_Color;
out vec4 v_FragData;

void main() {
  float v = 1.0 - abs(f_Edge);
  v = smoothstep(0.65, 0.7, v * f_Inner);
  v_FragData = mix(v_Color, vec4(0.0), v);
}` + "\x00"

const LINES_VERTEX = `#version 150
in vec2 v_Position;
in vec2 v_Normal;
in float f_Miter;
uniform float f_Thickness;
uniform mat4 m_ModelView;
uniform mat4 m_Projection;
out float f_Edge;

void main() {
    f_Edge = sign(f_Miter);

    //push the point along its normal by half thickness
    vec2 position = v_Position.xy + vec2(v_Normal * f_Thickness/2.0 * f_Miter);
    gl_Position = m_Projection * m_ModelView * vec4(position, 0.0, 1.0);
}` + "\x00"

func NewLinesRenderer(bounds, screen Rectangle) (lr *LinesRenderer, err error) {
	var (
		program uint32
		vbos    = make([]uint32, 2)
		r       *Renderer
		point   TexturedPoint
	)
	if program, err = BuildProgram(LINES_VERTEX, LINES_FRAGMENT); err != nil {
		return
	}
	if r, err = NewRenderer(bounds, screen); err != nil {
		return
	}
	gl.GenBuffers(2, &vbos[0])
	lr = &LinesRenderer{
		Renderer:      r,
		program:       program,
		buffer:        vbos[0],
		indexBuffer:   vbos[1],
		bufferBytes:   0,
		positionLoc:   uint32(gl.GetAttribLocation(program, gl.Str("v_Position\x00"))),
		normalLoc:     uint32(gl.GetAttribLocation(program, gl.Str("v_Normal\x00"))),
		miterLoc:      uint32(gl.GetAttribLocation(program, gl.Str("f_Miter\x00"))),
		modelviewLoc:  gl.GetUniformLocation(program, gl.Str("m_ModelView\x00")),
		projectionLoc: gl.GetUniformLocation(program, gl.Str("m_Projection\x00")),
		thicknessLoc:  gl.GetUniformLocation(program, gl.Str("f_Thickness\x00")),
		colorLoc:      gl.GetUniformLocation(program, gl.Str("v_Color\x00")),
		innerLoc:      gl.GetUniformLocation(program, gl.Str("f_Inner\x00")),
		offPosition:   gl.PtrOffset(int(unsafe.Offsetof(point.X))),
		offNormal:     gl.PtrOffset(int(unsafe.Offsetof(point.TextureX))),
		offMiter:      gl.PtrOffset(int(unsafe.Offsetof(point.Z))),
		stride:        int32(unsafe.Sizeof(point)),
	}
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (lr *LinesRenderer) Bind() (err error) {
	gl.UseProgram(lr.program)
	gl.BindBuffer(gl.ARRAY_BUFFER, lr.buffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, lr.indexBuffer)
	gl.EnableVertexAttribArray(lr.positionLoc)
	gl.EnableVertexAttribArray(lr.normalLoc)
	gl.EnableVertexAttribArray(lr.miterLoc)
	gl.VertexAttribPointer(lr.positionLoc, 2, gl.FLOAT, false, lr.stride, lr.offPosition)
	gl.VertexAttribPointer(lr.normalLoc, 2, gl.FLOAT, false, lr.stride, lr.offNormal)
	gl.VertexAttribPointer(lr.miterLoc, 1, gl.FLOAT, false, lr.stride, lr.offMiter)
	gl.UniformMatrix4fv(lr.projectionLoc, 1, false, &lr.Renderer.projection[0])
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (lr *LinesRenderer) Unbind() (err error) {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (lr *LinesRenderer) Draw(line *LineGeometry, mv mgl32.Mat4, style *LineStyle) (err error) {
	var (
		dataBytes    int   = len(line.Vertices) * int(lr.stride)
		indexBytes   int   = len(line.Indices) * int(lr.stride)
		elementCount int32 = int32(len(line.Indices))
		r, g, b, a         = style.Color.RGBA()
	)
	gl.Uniform1f(lr.thicknessLoc, style.Thickness)
	gl.Uniform1f(lr.innerLoc, style.Inner)
	gl.Uniform4f(lr.colorLoc, float32(r)/255.0, float32(g)/255.0, float32(b)/255.0, float32(a)/255.0)
	gl.UniformMatrix4fv(lr.modelviewLoc, 1, false, &mv[0])
	if dataBytes > lr.bufferBytes {
		lr.bufferBytes = dataBytes
		gl.BufferData(gl.ARRAY_BUFFER, dataBytes, gl.Ptr(line.Vertices), gl.STREAM_DRAW)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexBytes, gl.Ptr(line.Indices), gl.STREAM_DRAW)
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, dataBytes, gl.Ptr(line.Vertices))
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, indexBytes, gl.Ptr(line.Indices))
	}
	gl.DrawElements(gl.TRIANGLES, elementCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (lr *LinesRenderer) Delete() (err error) {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &lr.buffer)
	gl.DeleteBuffers(1, &lr.indexBuffer)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

type LineStyle struct {
	Thickness float32
	Color     color.Color
	Inner     float32
}
