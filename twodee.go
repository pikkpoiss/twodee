// Copyright 2012 Arne Roomann-Kurrik
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
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/banthar/gl"
	"github.com/jteeuwen/glfw"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type Window struct {
	Width  int
	Height int
}

func (w *Window) Opened() bool {
	return glfw.WindowParam(glfw.Opened) == 1
}

type System struct {
	Textures []gl.Texture
}

func Init() (sys *System, err error) {
	if err = glfw.Init(); err != nil {
		return
	}
	sys = &System{}
	sys.Textures = make([]gl.Texture, 0)
	return
}

func (s *System) Key(key int) int {
	return glfw.Key(key)
}

func (s *System) Terminate() {
	glfw.Terminate()
	s.Textures = nil
}

func (s *System) Open(win *Window) (err error) {
	err = glfw.OpenWindow(
		win.Width,
		win.Height,
		0, 0, 0, 0, 0, 0,
		glfw.Windowed)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(win.Width), float64(win.Height), 0, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.Enable(gl.TEXTURE_2D)
	_, err = s.LoadTexture("examples/basic/texture.png", IntNearest)
	return
}

func (s *System) LoadTexture(path string, smoothing int) (index int, err error) {
	var (
		file *os.File
		img  image.Image
		data *bytes.Buffer
	)
	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)

	if file, err = os.Open(path); err != nil {
		return
	}
	if img, err = png.Decode(file); err != nil {
		return
	}
	if data, err = EncodeTGA(" ", img); err != nil {
		return
	}
	if !glfw.LoadMemoryTexture2D(data.Bytes(), 0) {
		err = fmt.Errorf("Failed to load texture: %v", path)
		return
	}
	switch smoothing {
	default:
		fallthrough
	case IntNearest:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	case IntLinear:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}
	s.Textures = append(s.Textures, texture)
	index = len(s.Textures)
	return
}

func EncodeTGA(name string, img image.Image) (buf *bytes.Buffer, err error) {
	var (
		bounds image.Rectangle = img.Bounds()
		ident  []byte          = []byte(name)
		width  []byte          = make([]byte, 2)
		height []byte          = make([]byte, 2)
		nrgba  *image.NRGBA
		data   []byte
	)
	binary.LittleEndian.PutUint16(width, uint16(bounds.Dx()))
	binary.LittleEndian.PutUint16(height, uint16(bounds.Dy()))

	// See http://paulbourke.net/dataformats/tga/
	buf = &bytes.Buffer{}
	buf.WriteByte(byte(len(ident)))
	buf.WriteByte(0)
	buf.WriteByte(2) // uRGBI
	buf.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0})
	buf.Write([]byte(width))
	buf.Write([]byte(height))
	buf.WriteByte(32) // Bits per pixel
	buf.WriteByte(8)
	if buf.Len() != 18 {
		err = fmt.Errorf("TGA header is not 18 bytes: %v", buf.Len())
		return
	}

	nrgba = image.NewNRGBA(bounds)
	draw.Draw(nrgba, bounds, img, bounds.Min, draw.Src)
	buf.Write(ident)
	data = make([]byte, bounds.Dx()*bounds.Dy()*4)
	var (
		lineLength int = bounds.Dx() * 4
		destOffset int = len(data) - lineLength
	)
	for srcOffset := 0; srcOffset < len(nrgba.Pix); {
		var (
			dest   = data[destOffset : destOffset+lineLength]
			source = nrgba.Pix[srcOffset : srcOffset+nrgba.Stride]
		)
		copy(dest, source)
		destOffset -= lineLength
		srcOffset += nrgba.Stride
	}
	for x := 0; x < len(data); {
		buf.WriteByte(data[x+2])
		buf.WriteByte(data[x+1])
		buf.WriteByte(data[x+0])
		buf.WriteByte(data[x+3])
		x += 4
	}
	return
}

func (s *System) Paint() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	s.Textures[0].Unbind(gl.TEXTURE_2D)
	gl.Begin(gl.LINES)
	gl.Vertex2f(0, 0)
	gl.Vertex2f(640, 480)
	gl.End()

	s.Textures[0].Bind(gl.TEXTURE_2D)
	gl.Begin(gl.QUADS)
	gl.Vertex2f(50, 50)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(100, 50)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(100, 100)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(50, 100)
	gl.TexCoord2f(0, 1)
	gl.End()
	glfw.SwapBuffers()
}
