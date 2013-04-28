// Copyright 2013 Arne Roomann-Kurrik
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
	"github.com/kurrik/gltext"
	"os"
)

type Font struct {
	font *gltext.Font
}

func LoadFont(path string, scale int32) (font *Font, err error) {
	var (
		glf  *gltext.Font
		fd   *os.File
		low  rune = 32
		high rune = 127
		dir       = gltext.LeftToRight
	)
	if fd, err = os.Open(path); err != nil {
		return
	}
	defer fd.Close()

	if glf, err = gltext.LoadTruetype(fd, scale, low, high, dir); err != nil {
		return
	}

	font = &Font{
		font: glf,
	}
	return
}

func (f *Font) Printf(x float64, y float64, format string, a ...interface{}) (err error) {
	var (
		str string
		//sw  int
		//sh  int
	)
	str = fmt.Sprintf(format, a...)
	//sw, sh := f.font.Metrics(str)
	//gl.Color4f(1.0, 0.1, 0.1, 1)
	//gl.Rectd(x, y, x+float64(sw), y+float64(sh))
	gl.Color4f(1, 1, 1, 1)
	err = f.font.Printf(float32(x), float32(y), str)
	return
}
