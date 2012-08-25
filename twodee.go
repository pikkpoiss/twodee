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
	Title  string
}

func (w *Window) Opened() bool {
	return glfw.WindowParam(glfw.Opened) == 1
}

type Texture struct {
	texture gl.Texture
	Width   int
	Height  int
}

func LoadTexture(path string, smoothing int) (texture *Texture, err error) {
	var (
		file      *os.File
		img       image.Image
		bounds    image.Rectangle
		data      *bytes.Buffer
		gltexture gl.Texture
	)
	gltexture = gl.GenTexture()
	gltexture.Bind(gl.TEXTURE_2D)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()
	if img, err = png.Decode(file); err != nil {
		return
	}
	if data, err = EncodeTGA(path, img); err != nil {
		return
	}
	if !glfw.LoadMemoryTexture2D(data.Bytes(), 0) {
		err = fmt.Errorf("Failed to load texture: %v", path)
		return
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, smoothing)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, smoothing)
	bounds = img.Bounds()
	texture = &Texture{
		texture: gltexture,
		Width:   bounds.Dx(),
		Height:  bounds.Dy(),
	}
	return
}

func (t *Texture) Bind() {
	t.texture.Bind(gl.TEXTURE_2D)
}

func (t *Texture) Unbind() {
	t.texture.Unbind(gl.TEXTURE_2D)
}

type System struct {
	Textures map[string]*Texture
}

func Init() (sys *System, err error) {
	if err = glfw.Init(); err != nil {
		return
	}
	sys = &System{}
	sys.Textures = make(map[string]*Texture, 0)
	return
}

func (s *System) SetKeyCallback(handler glfw.KeyHandler) {
	glfw.SetKeyCallback(handler)
}

func (s *System) Key(key int) int {
	return glfw.Key(key)
}

func (s *System) SetCharCallback(handler glfw.CharHandler) {
	glfw.SetCharCallback(handler)
}

func (s *System) Terminate() {
	glfw.Terminate()
	s.Textures = nil
}

func (s *System) setProjection(win *Window) {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(win.Width), float64(win.Height), 0, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)
}

func (s *System) Open(win *Window) (err error) {
	glfw.SetWindowSizeCallback(func(w, h int) {
		win.Width = w
		win.Height = h
		s.setProjection(win)
	})
	err = glfw.OpenWindow(
		win.Width,
		win.Height,
		0, 0, 0, 0, 0, 0,
		glfw.Windowed)
	glfw.SetWindowTitle(win.Title)
	win.Width, win.Height = glfw.WindowSize()
	s.setProjection(win)
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	return
}

func (s *System) clamp(i int, max int) gl.GLclampf {
	return gl.GLclampf(float32(i) / float32(max))
}

func (s *System) SetClearColor(r int, g int, b int, a int) {
	gl.ClearColor(s.clamp(r, 255), s.clamp(g, 255), s.clamp(b, 255), s.clamp(a, 255))
}

func (s *System) LoadTexture(name string, path string, inter int) (err error) {
	var texture *Texture
	if texture, err = LoadTexture(path, inter); err != nil {
		return
	}
	s.Textures[name] = texture
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

func (s *System) Paint(scene *Scene) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	scene.Draw()
	glfw.SwapBuffers()
}
