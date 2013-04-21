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
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"image"
	"image/draw"
	"image/png"
	"os"
	"time"
)

type Window struct {
	Width      int
	Height     int
	Title      string
	View       Rectangle
	Fullscreen bool
	Scale      int
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
		trimbounds = image.Rect(0, 0, bounds.Dx(), bounds.Dy()-1)
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

func (t *Texture) Bind() {
	t.texture.Bind(gl.TEXTURE_2D)
}

func (t *Texture) Unbind() {
	t.texture.Unbind(gl.TEXTURE_2D)
}

func (t *Texture) Dispose() {
	t.texture.Delete()
}

type System struct {
	Textures      map[string]*Texture
	Framebuffer   *Framebuffer
	Overlay       *Framebuffer
	OverlayCamera *Camera
	Win           *Window
	LastPaint     time.Time
	resizeHandler glfw.WindowSizeHandler
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

func (s *System) SetSizeCallback(handler glfw.WindowSizeHandler) {
	s.resizeHandler = handler
}

func (s *System) Terminate() {
	for _, t := range s.Textures {
		t.Dispose()
	}
	s.Framebuffer.Dispose()
	s.Overlay.Dispose()
	glfw.Terminate()
}

func (s *System) resize() (err error) {
	s.Win.Width, s.Win.Height = glfw.WindowSize()
	if s.Framebuffer != nil {
		s.Framebuffer.Dispose()
	}
	if s.Overlay != nil {
		s.Overlay.Dispose()
	}
	var (
		fbw = float64(s.Win.Width) / float64(s.Win.Scale)
		fbh = float64(s.Win.Height) / float64(s.Win.Scale)
	)
	if s.Framebuffer, err = NewFramebuffer(int(fbw), int(fbh)); err != nil {
		return
	}
	if s.Overlay, err = NewFramebuffer(s.Win.Width, s.Win.Height); err != nil {
		return
	}
	s.OverlayCamera = NewCamera(0, 0, float64(s.Win.Width), float64(s.Win.Height))
	if s.resizeHandler != nil {
		s.resizeHandler(s.Win.Width, s.Win.Height)
	}
	return
}

func (s *System) Open(win *Window) (err error) {
	s.Win = win
	if win.Scale < 1 {
		win.Scale = 1
	}
	mode := glfw.Windowed
	if win.Fullscreen {
		mode = glfw.Fullscreen
	}
	if err = glfw.OpenWindow(win.Width, win.Height, 0, 0, 0, 0, 0, 0, mode); err != nil {
		return
	}
	gl.Init()
	v1, v2, v3 := glfw.GLVersion()
	fmt.Printf("OpenGL version: %v %v %v\n", v1, v2, v3)
	fmt.Printf("Framebuffer supported: %v\n", glfw.ExtensionSupported("GL_EXT_framebuffer_object"))
	glfw.SetSwapInterval(1) // Limit to refresh
	glfw.SetWindowTitle(win.Title)
	glfw.SetWindowSizeCallback(func(w, h int) {
		fmt.Printf("Resizing window to %v, %v\n", w, h)
		s.resize()
	})
	err = s.resize()
	return
}

func (s *System) clamp(i int, max int) gl.GLclampf {
	return gl.GLclampf(float64(i) / float64(max))
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

func (s *System) Paint(scene *Scene) {
	var (
		now time.Time
		fps float64
	)
	now = time.Now()
	if !s.LastPaint.IsZero() {
		fps = 1 / now.Sub(s.LastPaint).Seconds()
	}
	s.LastPaint = now

	s.Framebuffer.Bind()
	scene.Camera.SetProjection()
	gl.ClearColor(0.0, 0.0, 0.0, 0)
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	scene.Draw()
	gl.Flush()
	s.Framebuffer.Unbind()
	
	s.Overlay.Bind()
	s.OverlayCamera.SetProjection()
	gl.ClearColor(0.0, 0.0, 0.0, 0)
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	scene.Font.Printf(0, 0, "FPS %6.2f", fps)
	gl.Flush()
	s.Overlay.Unbind()

	gl.ClearColor(0.0, 0.0, 0.0, 0)
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	s.Framebuffer.Draw(s.Win.Width, s.Win.Height)
	s.Overlay.Draw(s.Win.Width, s.Win.Height)
	glfw.SwapBuffers()
}
