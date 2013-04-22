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
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"image"
)

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

type Framebuffer struct {
	Buffer  gl.Framebuffer
	Texture gl.Texture
	Width   int
	Height  int
}

func NewFramebuffer(w int, h int) (fb *Framebuffer, err error) {
	var (
		buffer  gl.Framebuffer
		texture gl.Texture
	)
	buffer = gl.GenFramebuffer()
	buffer.Bind()
	texture = gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture, 0)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0)
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		err = fmt.Errorf("Framebuffer could not be set up")
		return
	}
	fb = &Framebuffer{
		Buffer:  buffer,
		Texture: texture,
		Width:   w,
		Height:  h,
	}
	return
}

func (fb *Framebuffer) Dispose() {
	fb.Buffer.Delete()
	fb.Texture.Delete()
}

func (fb *Framebuffer) Bind() {
	fb.Buffer.Bind()
	gl.Viewport(0, 0, fb.Width, fb.Height)
	gl.Enable(gl.TEXTURE_2D)
	gl.Disable(gl.DEPTH_TEST)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
}

func (fb *Framebuffer) Unbind() {
	fb.Buffer.Unbind()
}

func (fb *Framebuffer) Draw(w int, h int) {
	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, 1, 0, 1, 1, -1)
	gl.Enable(gl.TEXTURE_2D)
	fb.Texture.Bind(gl.TEXTURE_2D)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Begin(gl.QUADS)
	gl.TexCoord2d(0, 0)
	gl.Vertex3f(0.0, 0.0, 0)
	gl.TexCoord2d(1, 0)
	gl.Vertex3f(1.0, 0.0, 0)
	gl.TexCoord2d(1, 1)
	gl.Vertex3f(1.0, 1.0, 0)
	gl.TexCoord2d(0, 1)
	gl.Vertex3f(0.0, 1.0, 0)
	gl.End()
	gl.Flush()
}
