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
	"time"
)

type System struct {
	Textures      map[string]*Texture
	Framebuffer   *Framebuffer
	Overlay       *Framebuffer
	OverlayCamera *Camera
	Win           *Window
	LastPaint     time.Time
	Font          *Font
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

type MouseMoveHandler func(x int, y int)

func (s *System) SetMouseMoveCallback(handler MouseMoveHandler) {
	glfw.SetMousePosCallback(glfw.MousePosHandler(handler))
}

type KeyHandler func(key int, state int)

func (s *System) SetKeyCallback(handler KeyHandler) {
	glfw.SetKeyCallback(glfw.KeyHandler(handler))
}

func (s *System) Key(key int) int {
	return glfw.Key(key)
}

type CharHandler func(key int, state int)

func (s *System) SetCharCallback(handler CharHandler) {
	glfw.SetCharCallback(glfw.CharHandler(handler))
}

type SizeHandler func(w int, h int)

func (s *System) SetSizeCallback(handler SizeHandler) {
	s.resizeHandler = glfw.WindowSizeHandler(handler)
}

type CloseHandler func() int

func (s *System) SetCloseCallback(handler CloseHandler) {
	glfw.SetWindowCloseCallback(glfw.WindowCloseHandler(handler))
}

type ScrollHandler func(pos int)

func (s *System) SetScrollCallback(handler ScrollHandler) {
	glfw.SetMouseWheelCallback(glfw.MouseWheelHandler(handler))
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
	var (
		oldframebuffer *Framebuffer = s.Framebuffer
		oldoverlay     *Framebuffer = s.Overlay
		fbw                         = float64(s.Win.Width) / float64(s.Win.Scale)
		fbh                         = float64(s.Win.Height) / float64(s.Win.Scale)
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
	if oldframebuffer != nil {
		oldframebuffer.Dispose()
	}
	if oldoverlay != nil {
		oldoverlay.Dispose()
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
	s.SetClearColor(0, 0, 0, 0)
	//glfw.SetSwapInterval(1) // Limit to refresh
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
	gl.ClearDepth(1.0)
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

func (s *System) SetFont(font *Font) {
	s.Font = font
}

func (s *System) Paint(scene Visible) {
	var (
		now time.Time
		fps float64
	)
	now = time.Now()
	if !s.LastPaint.IsZero() {
		fps = 1.0 / now.Sub(s.LastPaint).Seconds()
	}

	s.Framebuffer.Bind()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	scene.Draw()
	gl.Flush()
	s.Framebuffer.Unbind()

	if s.Font != nil {
		s.Overlay.Bind()
		s.OverlayCamera.SetProjection()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		s.Font.Printf(0, 0, "FPS %6.2f", fps)
		gl.Flush()
		s.Overlay.Unbind()
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	s.Framebuffer.Draw(s.Win.Width, s.Win.Height)
	s.Overlay.Draw(s.Win.Width, s.Win.Height)
	glfw.SwapBuffers()

	s.LastPaint = now
}
