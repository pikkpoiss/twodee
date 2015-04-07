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
	"unsafe"
)

type PointsRenderer struct {
	*Renderer
	program          uint32
	positionLoc      uint32
	translationLoc   uint32
	rotationLoc      uint32
	scaleLoc         uint32
	textureLoc       uint32
	textureUnitLoc   int32
	modelViewLoc     int32
	projectionLoc    int32
	vertexVBO        uint32
	instanceVBO      uint32
	vertexBytes      uint32
	instanceBytes    uint32
	sizePoint        uint32
	sizeAttr         uint32
	offPointX        unsafe.Pointer
	offPointTextureX unsafe.Pointer
	offAttrX         unsafe.Pointer
	offAttrRotationX unsafe.Pointer
	offAttrScaleX    unsafe.Pointer
}

const POINTS_FRAGMENT = `#version 150
precision mediump float;

uniform sampler2D TextureUnit;
in vec2 v_TextureCoordinates;
out vec4 v_FragData;

void main() {
  vec4 color = texture(TextureUnit, v_TextureCoordinates);
  //if (color.a < 0.1) {
  //  discard;
  //}
  v_FragData = color;
  //v_FragData = vec4(color.x, 1.0, 0.0, color.a);
}` + "\x00"

const POINTS_VERTEX = `#version 150

in vec4 v_Position;
in vec3 v_Translation;
in vec3 v_Rotation;
in vec3 v_Scale;
in vec2 TextureCoordinates;

uniform mat4 m_ProjectionMatrix;

out vec2 v_TextureCoordinates;

void main() {
  mat4 Translation = mat4(
    vec4(1.0, 0.0, 0.0, v_Translation.x),
    vec4(0.0, 1.0, 0.0, v_Translation.y),
    vec4(0.0, 0.0, 1.0, v_Translation.z),
    vec4(0.0, 0.0, 0.0,             1.0)
  );

  mat4 Scale = mat4(
    vec4(v_Scale.x,       0.0,       0.0, 0.0),
    vec4(      0.0, v_Scale.y,       0.0, 0.0),
    vec4(      0.0,       0.0, v_Scale.z, 0.0),
    vec4(      0.0,       0.0,       0.0, 1.0)
  );

  float rotCos = float(cos(v_Rotation.x));
  float rotSin = float(sin(v_Rotation.x));

  mat4 Rotation = mat4(
    vec4(rotCos, -rotSin, 0.0, 0.0),
    vec4(rotSin,  rotCos, 0.0, 0.0),
    vec4(   0.0,     0.0, 1.0, 0.0),
    vec4(   0.0,     0.0, 0.0, 1.0)
  );

  v_TextureCoordinates = TextureCoordinates;
  gl_Position =  m_ProjectionMatrix * (v_Position * Scale * Rotation * Translation);
}` + "\x00"

func NewPointsRenderer(bounds, screen Rectangle) (tr *PointsRenderer, err error) {
	var (
		program            uint32
		vbos               = make([]uint32, 2)
		r                  *Renderer
		texturedPoint      TexturedPoint
		instanceAttributes InstanceAttributes
	)
	if program, err = BuildProgram(POINTS_VERTEX, POINTS_FRAGMENT); err != nil {
		return
	}
	gl.GenBuffers(2, &vbos[0])
	if r, err = NewRenderer(bounds, screen); err != nil {
		return
	}
	tr = &PointsRenderer{
		Renderer:         r,
		program:          program,
		vertexVBO:        vbos[0],
		instanceVBO:      vbos[1],
		vertexBytes:      0,
		instanceBytes:    0,
		positionLoc:      uint32(gl.GetAttribLocation(program, gl.Str("v_Position\x00"))),
		translationLoc:   uint32(gl.GetAttribLocation(program, gl.Str("v_Translation\x00"))),
		rotationLoc:      uint32(gl.GetAttribLocation(program, gl.Str("v_Rotation\x00"))),
		scaleLoc:         uint32(gl.GetAttribLocation(program, gl.Str("v_Scale\x00"))),
		textureLoc:       uint32(gl.GetAttribLocation(program, gl.Str("TextureCoordinates\x00"))),
		textureUnitLoc:   gl.GetUniformLocation(program, gl.Str("TextureUnit\x00")),
		modelViewLoc:     gl.GetUniformLocation(program, gl.Str("ModelViewMatrix\x00")),
		projectionLoc:    gl.GetUniformLocation(program, gl.Str("m_ProjectionMatrix\x00")),
		sizePoint:        uint32(unsafe.Sizeof(texturedPoint)),
		sizeAttr:         uint32(unsafe.Sizeof(instanceAttributes)),
		offPointX:        gl.PtrOffset(int(unsafe.Offsetof(texturedPoint.X))),
		offPointTextureX: gl.PtrOffset(int(unsafe.Offsetof(texturedPoint.TextureX))),
		offAttrX:         gl.PtrOffset(int(unsafe.Offsetof(instanceAttributes.X))),
		offAttrRotationX: gl.PtrOffset(int(unsafe.Offsetof(instanceAttributes.RotationX))),
		offAttrScaleX:    gl.PtrOffset(int(unsafe.Offsetof(instanceAttributes.ScaleX))),
	}
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (tr *PointsRenderer) Draw(instances *InstanceList) error {
	var (
		bytesNeeded uint32
		stride      uint32
		count       uint32
		data        unsafe.Pointer
	)
	gl.UseProgram(tr.program)
	gl.Uniform1i(tr.textureUnitLoc, 0)

	// Vertex data binding.
	gl.BindBuffer(gl.ARRAY_BUFFER, tr.vertexVBO)
	stride = tr.sizePoint
	count = uint32(len(instances.Geometry))
	bytesNeeded = stride * count
	data = gl.Ptr(instances.Geometry)
	if bytesNeeded > tr.vertexBytes {
		gl.BufferData(gl.ARRAY_BUFFER, int(bytesNeeded), data, gl.STREAM_DRAW)
		tr.vertexBytes = bytesNeeded
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, int(bytesNeeded), data)
	}
	gl.EnableVertexAttribArray(tr.positionLoc)
	gl.EnableVertexAttribArray(tr.textureLoc)
	gl.VertexAttribPointer(tr.positionLoc, 3, gl.FLOAT, false, int32(stride), tr.offPointX)
	gl.VertexAttribPointer(tr.textureLoc, 2, gl.FLOAT, false, int32(stride), tr.offPointTextureX)

	// Instance data binding
	gl.BindBuffer(gl.ARRAY_BUFFER, tr.instanceVBO)
	stride = tr.sizeAttr
	count = uint32(len(instances.Instances))
	bytesNeeded = stride * count
	data = gl.Ptr(instances.Instances)
	if bytesNeeded > tr.instanceBytes {
		gl.BufferData(gl.ARRAY_BUFFER, int(bytesNeeded), data, gl.STREAM_DRAW)
		tr.instanceBytes = bytesNeeded
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, int(bytesNeeded), data)
	}
	gl.EnableVertexAttribArray(tr.translationLoc)
	gl.EnableVertexAttribArray(tr.rotationLoc)
	gl.EnableVertexAttribArray(tr.scaleLoc)
	gl.VertexAttribPointer(tr.translationLoc, 3, gl.FLOAT, false, int32(stride), tr.offAttrX)
	gl.VertexAttribPointer(tr.rotationLoc, 3, gl.FLOAT, false, int32(stride), tr.offAttrRotationX)
	gl.VertexAttribPointer(tr.scaleLoc, 3, gl.FLOAT, false, int32(stride), tr.offAttrScaleX)
	gl.VertexAttribDivisor(tr.translationLoc, 1)
	gl.VertexAttribDivisor(tr.rotationLoc, 1)
	gl.VertexAttribDivisor(tr.scaleLoc, 1)

	// Projection
	gl.UniformMatrix4fv(tr.projectionLoc, 1, false, &tr.Renderer.projection[0])

	// Actually draw.
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(len(instances.Geometry)), int32(len(instances.Instances)))

	// Undo instance attr repetition.
	gl.VertexAttribDivisor(tr.translationLoc, 0)
	gl.VertexAttribDivisor(tr.rotationLoc, 0)
	gl.VertexAttribDivisor(tr.scaleLoc, 0)

	return nil
}

func (tr *PointsRenderer) Delete() error {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &tr.instanceVBO)
	gl.DeleteBuffers(1, &tr.vertexVBO)
	return nil
}
