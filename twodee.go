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

type Point struct {
	X float32
	Y float32
}

func Pt(x float32, y float32) Point {
	return Point{X: x, Y: y}
}

type Rectangle struct {
	Min Point
	Max Point
}

func Rect(x1 float32, y1 float32, x2 float32, y2 float32) Rectangle {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	return Rectangle{Min:Pt(x1, y1), Max:Pt(x2, y2)}
}

func (r Rectangle) Overlaps(s Rectangle) bool {
	return r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
	       r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

type Window struct {
	Width  int
	Height int
	Title  string
	View   Rectangle
}

func (w *Window) Opened() bool {
	return glfw.WindowParam(glfw.Opened) == 1
}

type Texture struct {
	texture gl.Texture
	Width   int
	Height  int
	Frames  [][]int
}

func LoadPNG(path string) (img image.Image, err error) {
	var file *os.File
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()
	img, err = png.Decode(file)
	return
}

func LoadTexture(path string, smoothing int, framewidth int) (texture *Texture, err error) {
	var (
		img    image.Image
		bounds image.Rectangle
		gltex  gl.Texture
	)
	if img, err = LoadPNG(path); err != nil {
		return
	}
	bounds = img.Bounds()
	if gltex, err = GetGLTexture(img, smoothing); err != nil {
		return
	}
	texture = &Texture{
		texture: gltex,
		Width:   bounds.Dx(),
		Height:  bounds.Dy(),
		Frames:  make([][]int, 0),
	}
	frames := bounds.Dx() / framewidth
	for i := 0; i < frames; i++ {
		texture.Frames = append(texture.Frames, []int{
			i * framewidth,
			(i + 1) * framewidth,
		})
	}
	return
}

func LoadVarWidthTexture(path string, smoothing int) (texture *Texture, err error) {
	var (
		img   image.Image
		trim  *image.NRGBA
		gltex gl.Texture
	)
	if img, err = LoadPNG(path); err != nil {
		return
	}
	var (
		bounds     = img.Bounds()
		trimbounds = image.Rect(0, 0, bounds.Dx(), bounds.Dy() - 1)
		trimpoint  = image.Pt(0, 1)
	)
	trim = image.NewNRGBA(trimbounds)
	draw.Draw(trim, trimbounds, img, trimpoint, draw.Src)
	if gltex, err = GetGLTexture(trim, smoothing); err != nil {
		return
	}
	texture = &Texture{
		texture: gltex,
		Width:   trimbounds.Dx(),
		Height:  trimbounds.Dy(),
		Frames:  make([][]int, 0),
	}
	var (
		aprime uint32 = 0
		pair          = make([]int, 2)
		x             = 0
	)
	for ; x < bounds.Dx(); x++ {
		_, _, _, a := img.At(x, 0).RGBA()
		if aprime == 0 && a > 0 {
			pair[0] = x
		} else if aprime > 0 && a == 0 {
			pair[1] = x
			texture.Frames = append(texture.Frames, pair)
			pair = make([]int, 2)
		}
		aprime = a
	}
	if pair[0] != 0 {
		pair[1] = x
		texture.Frames = append(texture.Frames, pair)
	}
	return
}

func GetGLTexture(img image.Image, smoothing int) (gltexture gl.Texture, err error) {
	var data *bytes.Buffer
	if data, err = EncodeTGA("texture", img); err != nil {
		return
	}
	gltexture = gl.GenTexture()
	gltexture.Bind(gl.TEXTURE_2D)
	if !glfw.LoadMemoryTexture2D(data.Bytes(), 0) {
		err = fmt.Errorf("Failed to load texture")
		return
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, smoothing)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, smoothing)
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
	win.View = Rect(0, 0, float32(win.Width) / 2, float32(win.Height) / 2)
	gl.Ortho(0, float64(win.View.Max.X), float64(win.View.Max.Y), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
}

func (s *System) Open(win *Window) (err error) {
	glfw.SetWindowSizeCallback(func(w, h int) {
		fmt.Printf("Resizing window to %v, %v\n", w, h)
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
	gl.Enable(gl.DEPTH_TEST)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	return
}

func (s *System) clamp(i int, max int) gl.GLclampf {
	return gl.GLclampf(float32(i) / float32(max))
}

func (s *System) SetClearColor(r int, g int, b int, a int) {
	gl.ClearColor(s.clamp(r, 255), s.clamp(g, 255), s.clamp(b, 255), s.clamp(a, 255))
}

func (s *System) LoadTexture(name string, path string, inter int, width int) (err error) {
	var texture *Texture
	if width > 0 {
		if texture, err = LoadTexture(path, inter, width); err != nil {
			return
		}
	} else {
		if texture, err = LoadVarWidthTexture(path, inter); err != nil {
			return
		}
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
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	scene.Draw()
	glfw.SwapBuffers()
}
