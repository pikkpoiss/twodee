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

type SoundEffect struct {
	chunk *mixer.Chunk
}

func NewSoundEffect(path string) (s *SoundEffect, err error) {
	s = &SoundEffect{
		chunk: mixer.LoadWAV(path),
	}
	if s.chunk == nil {
		err = fmt.Errorf("Could not load sound effect: %v", sdl.GetError())
	}
	return
}

func (s *SoundEffect) Delete() {
	s.chunk.Free()
}

func (s *SoundEffect) Play(times int) {
	s.chunk.PlayChannel(-1, times-1)
}

func (s *SoundEffect) PlayChannel(channel int, times int) {
	s.chunk.PlayChannel(channel, times-1)
}

func (s *SoundEffect) SetVolume(volume int) {
	s.chunk.Volume(volume)
}

func (s *SoundEffect) IsPlaying(channel int) int {
	return s.chunk.IsPlaying(channel)
}
