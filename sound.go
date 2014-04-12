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
	"github.com/kurrik/Go-SDL/mixer"
	"github.com/kurrik/Go-SDL/sdl"
)

func initSound() error {
	if code := sdl.Init(sdl.INIT_AUDIO); code == -1 {
		return fmt.Errorf("Could not init audio: %v", sdl.GetError())
	}
	var (
		freq     int    = mixer.DEFAULT_FREQUENCY
		format   uint16 = mixer.DEFAULT_FORMAT
		channels int    = mixer.DEFAULT_CHANNELS
	)
	if mixer.OpenAudio(freq, format, channels, 4096) != 0 {
		return fmt.Errorf("Could not init mixer: %v", sdl.GetError())
	}
	return nil
}

func cleanupSound() {
	mixer.CloseAudio()
	sdl.QuitSubSystem(sdl.INIT_AUDIO)
}

type Audio struct {
	music *mixer.Music
}

func NewAudio(path string) *Audio {
	return &Audio{
		music: mixer.LoadMUS(path),
	}
}

func (a *Audio) Delete() {
	a.music.Free()
}

func (a *Audio) Play(times int) {
	a.music.PlayMusic(times)
}
