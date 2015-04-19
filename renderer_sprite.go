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

type SpriteRenderer struct {
	*Renderer
	program           uint32
	translationLoc    uint32
	rotationLoc       uint32
	scaleLoc          uint32
	pointAdjLoc       uint32
	textureAdjLoc     uint32
	colorLoc          uint32
	textureUnitLoc    int32
	projectionLoc     int32
	instanceVBO       uint32
	instanceBytes     int
	offAttrX          unsafe.Pointer
	offAttrRotationX  unsafe.Pointer
	offAttrScaleX     unsafe.Pointer
	offAttrColor      unsafe.Pointer
	offAttrPointAdj   int
	offAttrTextureAdj int
}

const SPRITE_FRAGMENT = `#version 150
precision mediump float;

uniform sampler2D TextureUnit;
in vec2 v_TextureCoordinates;
in vec4 v_BaseColor;
out vec4 v_FragData;

void main() {
  vec4 color = texture(TextureUnit, v_TextureCoordinates);
  //if (color.a < 0.1) {
  //  discard;
  //}
  if (v_BaseColor.a > 0.1 && color.a < 0.1) {
    v_FragData = v_BaseColor;
  } else {
    v_FragData = color;
  }
}` + "\x00"

const SPRITE_VERTEX = `#version 150

in vec3 v_Translation;
in vec3 v_Rotation;
in vec3 v_Scale;
in vec4 v_Color;

in mat4 m_PointAdjustment;
in mat4 m_TextureAdjustment;

uniform mat4 m_ProjectionMatrix;

out vec4 v_BaseColor;
out vec2 v_TextureCoordinates;

const vec2 Points[] = vec2[6](
  vec2(-0.5, -0.5),
  vec2( 0.5,  0.5),
  vec2(-0.5,  0.5),
  vec2(-0.5, -0.5),
  vec2( 0.5, -0.5),
  vec2( 0.5,  0.5)
);

const vec2 TexturePoints[] = vec2[6](
  vec2(0.0, 0.0),
  vec2(1.0, 1.0),
  vec2(0.0, 1.0),
  vec2(0.0, 0.0),
  vec2(1.0, 0.0),
  vec2(1.0, 1.0)
);

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

  v_BaseColor = v_Color;
  v_TextureCoordinates = (vec4(TexturePoints[gl_VertexID], 0.0, 1.0) *
     m_TextureAdjustment).xy;

  gl_Position =  m_ProjectionMatrix *
      (vec4(Points[gl_VertexID], 0.0, 1.0) *
      m_PointAdjustment *
      Scale *
      Rotation *
      Translation);
}` + "\x00"

func NewSpriteRenderer(camera *Camera) (tr *SpriteRenderer, err error) {
	var (
		program     uint32
		vbos        = make([]uint32, 1)
		sprite      SpriteConfig
		frameOffset = int(unsafe.Offsetof(sprite.Frame))
		viewOffset  = int(unsafe.Offsetof(sprite.View))
	)
	if program, err = BuildProgram(SPRITE_VERTEX, SPRITE_FRAGMENT); err != nil {
		return
	}
	gl.GenBuffers(1, &vbos[0])
	tr = &SpriteRenderer{
		Renderer:          NewRenderer(camera),
		program:           program,
		instanceVBO:       vbos[0],
		instanceBytes:     0,
		translationLoc:    uint32(gl.GetAttribLocation(program, gl.Str("v_Translation\x00"))),
		rotationLoc:       uint32(gl.GetAttribLocation(program, gl.Str("v_Rotation\x00"))),
		scaleLoc:          uint32(gl.GetAttribLocation(program, gl.Str("v_Scale\x00"))),
		pointAdjLoc:       uint32(gl.GetAttribLocation(program, gl.Str("m_PointAdjustment\x00"))),
		textureAdjLoc:     uint32(gl.GetAttribLocation(program, gl.Str("m_TextureAdjustment\x00"))),
		colorLoc:          uint32(gl.GetAttribLocation(program, gl.Str("v_Color\x00"))),
		textureUnitLoc:    gl.GetUniformLocation(program, gl.Str("TextureUnit\x00")),
		projectionLoc:     gl.GetUniformLocation(program, gl.Str("m_ProjectionMatrix\x00")),
		offAttrX:          gl.PtrOffset(viewOffset + int(unsafe.Offsetof(sprite.View.X))),
		offAttrRotationX:  gl.PtrOffset(viewOffset + int(unsafe.Offsetof(sprite.View.RotationX))),
		offAttrColor:          gl.PtrOffset(int(unsafe.Offsetof(sprite.Color))),
		offAttrScaleX:     gl.PtrOffset(viewOffset + int(unsafe.Offsetof(sprite.View.ScaleX))),
		offAttrPointAdj:   frameOffset + int(unsafe.Offsetof(sprite.Frame.PointAdjustment)),
		offAttrTextureAdj: frameOffset + int(unsafe.Offsetof(sprite.Frame.TextureAdjustment)),
	}
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (tr *SpriteRenderer) Draw(instances []SpriteConfig) error {
	var (
		bytesNeeded int
		byteoffset  int
		count       int32
		data        unsafe.Pointer
		float       float32 = 0
		floatSize   uint32
		i           uint32
		offset      unsafe.Pointer
		sprite      SpriteConfig
		stride      int32
	)

	floatSize = uint32(unsafe.Sizeof(float))
	stride = int32(unsafe.Sizeof(sprite))

	gl.UseProgram(tr.program)
	gl.Uniform1i(tr.textureUnitLoc, 0)

	// Instance data binding
	gl.BindBuffer(gl.ARRAY_BUFFER, tr.instanceVBO)
	count = int32(len(instances))
	bytesNeeded = int(stride * count)
	data = gl.Ptr(instances)
	if bytesNeeded > tr.instanceBytes {
		gl.BufferData(gl.ARRAY_BUFFER, bytesNeeded, data, gl.STREAM_DRAW)
		tr.instanceBytes = bytesNeeded
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, bytesNeeded, data)
	}

	gl.EnableVertexAttribArray(tr.translationLoc)
	gl.VertexAttribPointer(tr.translationLoc, 3, gl.FLOAT, false, stride, tr.offAttrX)
	gl.VertexAttribDivisor(tr.translationLoc, 1)

	gl.EnableVertexAttribArray(tr.rotationLoc)
	gl.VertexAttribPointer(tr.rotationLoc, 3, gl.FLOAT, false, stride, tr.offAttrRotationX)
	gl.VertexAttribDivisor(tr.rotationLoc, 1)

	gl.EnableVertexAttribArray(tr.scaleLoc)
	gl.VertexAttribPointer(tr.scaleLoc, 3, gl.FLOAT, false, stride, tr.offAttrScaleX)
	gl.VertexAttribDivisor(tr.scaleLoc, 1)

	gl.EnableVertexAttribArray(tr.colorLoc)
	gl.VertexAttribPointer(tr.colorLoc, 4, gl.FLOAT, false, stride, tr.offAttrColor)
	gl.VertexAttribDivisor(tr.colorLoc, 1)

	for i = 0; i < 4; i++ {
		byteoffset = int(i * 4 * floatSize)
		offset = gl.PtrOffset(tr.offAttrPointAdj + byteoffset)
		gl.EnableVertexAttribArray(tr.pointAdjLoc + i)
		gl.VertexAttribPointer(tr.pointAdjLoc+i, 4, gl.FLOAT, false, stride, offset)
		gl.VertexAttribDivisor(tr.pointAdjLoc+i, 1)

		offset = gl.PtrOffset(tr.offAttrTextureAdj + byteoffset)
		gl.EnableVertexAttribArray(tr.textureAdjLoc + i)
		gl.VertexAttribPointer(tr.textureAdjLoc+i, 4, gl.FLOAT, false, stride, offset)
		gl.VertexAttribDivisor(tr.textureAdjLoc+i, 1)
	}
	// Projection
	gl.UniformMatrix4fv(tr.projectionLoc, 1, false, &tr.Renderer.Camera.Projection[0])

	// Actually draw.
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, int32(len(instances)))

	// Undo instance attr repetition.
	gl.VertexAttribDivisor(tr.translationLoc, 0)
	gl.VertexAttribDivisor(tr.rotationLoc, 0)
	gl.VertexAttribDivisor(tr.scaleLoc, 0)
	gl.VertexAttribDivisor(tr.colorLoc, 0)
	for i = 0; i < 4; i++ {
		gl.VertexAttribDivisor(tr.pointAdjLoc+i, 0)
		gl.VertexAttribDivisor(tr.textureAdjLoc+i, 0)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return nil
}

func (tr *SpriteRenderer) Delete() error {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &tr.instanceVBO)
	return nil
}
