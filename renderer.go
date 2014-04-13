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

type Renderer struct {
	worldBounds  Rectangle
	screenBounds Rectangle
	projection   *Matrix4
	inverse      *Matrix4
}

func NewRenderer(world, screen Rectangle) (r *Renderer, err error) {
	r = &Renderer{}
	r.SetScreenBounds(screen)
	err = r.SetWorldBounds(world)
	return
}

func (r *Renderer) SetScreenBounds(bounds Rectangle) {
	r.screenBounds = bounds
}

func (r *Renderer) SetWorldBounds(bounds Rectangle) (err error) {
	r.worldBounds = bounds
	r.projection = GetOrthoMatrix(bounds.Min.X, bounds.Max.X, bounds.Min.Y, bounds.Max.Y, 1, 0)
	r.inverse, err = GetInverseMatrix(r.projection)
	return
}

func (r *Renderer) ScreenToWorldCoords(x, y float32) (wx, wy float32) {
	// http://stackoverflow.com/questions/7692988/
	var (
		halfw = r.screenBounds.Max.X / 2.0
		halfh = r.screenBounds.Max.Y / 2.0
		xpct  = (x - halfw) / halfw
		ypct  = (halfh - y) / halfh
	)
	return Unproject(r.inverse, xpct, ypct)
}

func CreateVAO() (array gl.VertexArray, err error) {
	array = gl.GenVertexArray()
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR gl.GenVertexArray %X", e)
		return
	}
	array.Bind()
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR array.Bind %X", e)
		return
	}
	return
}

func CreateVBO(size int, data interface{}, usage gl.GLenum) (buffer gl.Buffer, err error) {
	buffer = gl.GenBuffer()
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR gl.GenBuffer %X", e)
		return
	}
	buffer.Bind(gl.ARRAY_BUFFER)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR buffer.Bind %X", e)
		return
	}
	gl.BufferData(gl.ARRAY_BUFFER, size, data, usage)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR gl.BufferData %X", e)
		return
	}
	buffer.Unbind(gl.ARRAY_BUFFER)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR buffer.Unbind %X", e)
		return
	}
	return
}

func CompileShader(stype gl.GLenum, source string) (shader gl.Shader, err error) {
	shader = gl.CreateShader(stype)
	shader.Source(source)
	shader.Compile()
	if status := shader.Get(gl.COMPILE_STATUS); status == 0 {
		err = fmt.Errorf("ERROR shader compile:\n%s", shader.GetInfoLog())
	}
	return
}

func LinkProgram(vertex gl.Shader, fragment gl.Shader) (program gl.Program, err error) {
	program = gl.CreateProgram()
	program.AttachShader(vertex)
	program.AttachShader(fragment)
	program.BindFragDataLocation(0, "v_FragData")
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR program.BindFragDataLocation %X", e)
		return
	}
	program.Link()
	if status := program.Get(gl.LINK_STATUS); status == 0 {
		err = fmt.Errorf("ERROR program link:\n%s", program.GetInfoLog())
	}
	return
}

func BuildProgram(vsrc string, fsrc string) (program gl.Program, err error) {
	var (
		vertex   gl.Shader
		fragment gl.Shader
	)
	if vertex, err = CompileShader(gl.VERTEX_SHADER, vsrc); err != nil {
		return
	}
	if fragment, err = CompileShader(gl.FRAGMENT_SHADER, fsrc); err != nil {
		return
	}
	return LinkProgram(vertex, fragment)
}

type TextRenderer struct {
	*Renderer
	VBO            gl.Buffer
	Program        gl.Program
	Texture        gl.Texture
	PositionLoc    gl.AttribLocation
	TextureLoc     gl.AttribLocation
	ScaleLoc       gl.UniformLocation
	TransLoc       gl.UniformLocation
	ProjectionLoc  gl.UniformLocation
	TextureUnitLoc gl.UniformLocation
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
}`

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
}`

func NewTextRenderer(screen Rectangle) (tr *TextRenderer, err error) {
	var (
		rect    []float32
		program gl.Program
		vbo     gl.Buffer
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
		PositionLoc:    program.GetAttribLocation("a_Position"),
		TextureLoc:     program.GetAttribLocation("a_TextureCoordinates"),
		TextureUnitLoc: program.GetUniformLocation("u_TextureUnit"),
		TransLoc:       program.GetUniformLocation("v_Trans"),
		ScaleLoc:       program.GetUniformLocation("v_Scale"),
		ProjectionLoc:  program.GetUniformLocation("m_ProjectionMatrix"),
		Width:          screen.Max.X - screen.Min.X,
		Height:         screen.Max.Y - screen.Min.Y,
	}
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (tr *TextRenderer) Bind() error {
	tr.Program.Use()
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

func (tr *TextRenderer) Draw(tex *Texture, x, y float32) (err error) {
	gl.ActiveTexture(gl.TEXTURE0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tex.Texture.Bind(gl.TEXTURE_2D)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.ScaleLoc.Uniform3f(float32(tex.Width), float32(tex.Height), 1)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.TransLoc.Uniform3f(x, y, 0)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	tr.VBO.Unbind(gl.ARRAY_BUFFER)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (tr *TextRenderer) Unbind() error {
	tr.VBO.Unbind(gl.ARRAY_BUFFER)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR: %X", e)
	}
	return nil
}

func (tr *TextRenderer) Delete() error {
	tr.VBO.Delete()
	return nil
}

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
