package twodee

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const GLOW_FRAGMENT = `#version 150
precision mediump float;

uniform sampler2D u_TextureUnit;
in vec2 v_TextureCoordinates;
uniform int Orientation;
uniform int BlurAmount;
uniform float BlurScale;
uniform float BlurStrength;
uniform vec2 BufferDimensions;
out vec4 v_FragData;
vec2 TexelSize = vec2(1.0 / BufferDimensions.x, 1.0 / BufferDimensions.y);

float Gaussian (float x, float deviation) {
  return (1.0 / sqrt(2.0 * 3.141592 * deviation)) * exp(-((x * x) / (2.0 * deviation)));
}


void main() {
  // Locals
  float halfBlur = float(BlurAmount) * 0.5;
  vec4 colour = vec4(0.0);
  vec4 texColour = vec4(0.0);

  // Gaussian deviation
  float deviation = halfBlur * 0.35;
  deviation *= deviation;
  float strength = 1.0 - BlurStrength;

  if ( Orientation == 0 ) {
    // Horizontal blur
    for (int i = 0; i < 10; ++i) {
      if ( i >= BlurAmount ) {
        break;
      }
      float offset = float(i) - halfBlur;
      texColour = texture(
        u_TextureUnit,
        v_TextureCoordinates + vec2(offset * TexelSize.x * BlurScale, 0.0)) * Gaussian(offset * strength, deviation);
      colour += texColour;
    }
  } else {
    // Vertical blur
    for (int i = 0; i < 10; ++i) {
      if ( i >= BlurAmount ) {
        break;
      }
      float offset = float(i) - halfBlur;
      texColour = texture(
        u_TextureUnit,
        v_TextureCoordinates + vec2(0.0, offset * TexelSize.y * BlurScale)) * Gaussian(offset * strength, deviation);
      colour += texColour;
    }
  }
  // Apply colour
  v_FragData = clamp(colour, 0.0, 1.0);
  v_FragData.w = 1.0;
}` + "\x00"

const GLOW_VERTEX = `#version 150

in vec4 a_Position;
in vec2 a_TextureCoordinates;

out vec2 v_TextureCoordinates;

void main() {
    v_TextureCoordinates = a_TextureCoordinates;
    gl_Position = a_Position;
}` + "\x00"

type GlowRenderer struct {
	GlowFb              uint32
	GlowTex             uint32
	BlurFb              uint32
	BlurTex             uint32
	shader              uint32
	positionLoc         uint32
	textureLoc          uint32
	orientationLoc      int32
	blurAmountLoc       int32
	blurScaleLoc        int32
	blurStrengthLoc     int32
	bufferDimensionsLoc int32
	textureUnitLoc      int32
	coords              uint32
	width               int
	height              int
	oldwidth            int
	oldheight           int
	blur                int
	strength            float32
	scale               float32
}

func NewGlowRenderer(w, h int, blur int, strength float32, scale float32) (r *GlowRenderer, err error) {
	r = &GlowRenderer{
		width:    w,
		height:   h,
		blur:     blur,
		strength: strength,
		scale:    scale,
	}
	_, _, r.oldwidth, r.oldheight = GetInteger4(gl.VIEWPORT)
	if r.shader, err = BuildProgram(GLOW_VERTEX, GLOW_FRAGMENT); err != nil {
		return
	}
	r.orientationLoc = gl.GetUniformLocation(r.shader, gl.Str("Orientation\x00"))
	r.blurAmountLoc = gl.GetUniformLocation(r.shader, gl.Str("BlurAmount\x00"))
	r.blurScaleLoc = gl.GetUniformLocation(r.shader, gl.Str("BlurScale\x00"))
	r.blurStrengthLoc = gl.GetUniformLocation(r.shader, gl.Str("BlurStrength\x00"))
	r.bufferDimensionsLoc = gl.GetUniformLocation(r.shader, gl.Str("BufferDimensions\x00"))
	r.positionLoc = uint32(gl.GetAttribLocation(r.shader, gl.Str("a_Position\x00")))
	r.textureLoc = uint32(gl.GetAttribLocation(r.shader, gl.Str("a_TextureCoordinates\x00")))
	r.textureUnitLoc = gl.GetUniformLocation(r.shader, gl.Str("u_TextureUnit\x00"))
	gl.BindFragDataLocation(r.shader, 0, gl.Str("v_FragData\x00"))
	var size float32 = 1.0
	var rect = []float32{
		-size, -size, 0.0, 0, 0,
		-size, size, 0.0, 0, 1,
		size, -size, 0.0, 1, 0,
		size, size, 0.0, 1, 1,
	}
	if r.coords, err = CreateVBO(len(rect)*4, rect, gl.STATIC_DRAW); err != nil {
		return
	}

	if r.GlowFb, r.GlowTex, err = r.initFramebuffer(w, h); err != nil {
		return
	}
	if r.BlurFb, r.BlurTex, err = r.initFramebuffer(w, h); err != nil {
		return
	}
	return
}

func (r *GlowRenderer) initFramebuffer(w, h int) (fb uint32, tex uint32, err error) {
	gl.GenFramebuffers(1, &fb)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb)

	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	gl.FramebufferTexture2D(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, tex, 0)
	if err = r.GetError(); err != nil {
		return
	}
	buffers := []uint32{gl.COLOR_ATTACHMENT0}
	gl.DrawBuffers(1, &buffers[0])

	var rb uint32
	gl.GenRenderbuffers(1, &rb)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rb)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.STENCIL_INDEX8, int32(w), int32(h))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.STENCIL_ATTACHMENT, gl.RENDERBUFFER, rb)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	return
}

func (r *GlowRenderer) SetStrength(s float32) {
	r.strength = s
}

func (r *GlowRenderer) GetError() error {
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("OpenGL error: %X", e)
	}
	var status = gl.CheckFramebufferStatus(gl.DRAW_FRAMEBUFFER)
	switch status {
	case gl.FRAMEBUFFER_COMPLETE:
		return nil
	case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
		return fmt.Errorf("Attachment point unconnected")
	case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
		return fmt.Errorf("Missing attachment")
	case gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER:
		return fmt.Errorf("Draw buffer")
	case gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER:
		return fmt.Errorf("Read buffer")
	case gl.FRAMEBUFFER_UNSUPPORTED:
		return fmt.Errorf("Unsupported config")
	default:
		return fmt.Errorf("Unknown framebuffer error: %X", status)
	}
}

func (r *GlowRenderer) Delete() error {
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteFramebuffers(1, &r.GlowFb)
	gl.DeleteTextures(1, &r.GlowTex)
	gl.DeleteFramebuffers(1, &r.BlurFb)
	gl.DeleteTextures(1, &r.BlurTex)
	gl.DeleteBuffers(1, &r.coords)
	return r.GetError()
}

func (r *GlowRenderer) Bind() error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.GlowFb)
	gl.Enable(gl.STENCIL_TEST)
	gl.Viewport(0, 0, int32(r.width), int32(r.height))
	gl.ClearStencil(0)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.StencilMask(0xFF) // Write to buffer
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	gl.StencilMask(0x00) // Don't write to buffer
	return nil
}

func (r *GlowRenderer) Draw() (err error) {
	gl.UseProgram(r.shader)
	gl.Uniform1i(r.textureUnitLoc, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.coords)
	gl.EnableVertexAttribArray(r.positionLoc)
	gl.VertexAttribPointer(r.positionLoc, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(r.textureLoc)
	gl.VertexAttribPointer(r.textureLoc, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.Uniform1i(r.blurAmountLoc, int32(r.blur))
	gl.Uniform1f(r.blurScaleLoc, r.scale)
	gl.Uniform1f(r.blurStrengthLoc, r.strength)
	gl.Uniform2f(r.bufferDimensionsLoc, float32(r.width), float32(r.height))

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.BlurFb)
	gl.Viewport(0, 0, int32(r.width), int32(r.height))
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, r.GlowTex)
	gl.Uniform1i(r.orientationLoc, 0)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	gl.Viewport(0, 0, int32(r.oldwidth), int32(r.oldheight))
	gl.BlendFunc(gl.ONE, gl.ONE)
	gl.BindTexture(gl.TEXTURE_2D, r.BlurTex)
	gl.Uniform1i(r.orientationLoc, 1)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	return nil
}

func (r *GlowRenderer) Unbind() error {
	gl.Viewport(0, 0, int32(r.oldwidth), int32(r.oldheight))
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Disable(gl.STENCIL_TEST)
	return nil
}

func (r *GlowRenderer) DisableOutput() {
	gl.ColorMask(false, false, false, false)
	gl.StencilFunc(gl.NEVER, 1, 0xFF)                // Never pass
	gl.StencilOp(gl.REPLACE, gl.REPLACE, gl.REPLACE) // Replace to ref=1
	gl.StencilMask(0xFF)                             // Write to buffer
}

func (r *GlowRenderer) EnableOutput() {
	gl.ColorMask(true, true, true, true)
	gl.StencilMask(0x00)              // No more writing
	gl.StencilFunc(gl.EQUAL, 0, 0xFF) // Only pass where stencil is 0
}

// Convenience function for glGetIntegerv
func GetInteger4(pname uint32) (v0, v1, v2, v3 int) {
	var values = []int32{0, 0, 0, 0}
	gl.GetIntegerv(pname, &values[0])
	v0 = int(values[0])
	v1 = int(values[1])
	v2 = int(values[2])
	v3 = int(values[3])
	return
}
