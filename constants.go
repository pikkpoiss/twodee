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

// Keyboard key definitions: 8-bit ISO-8859-1 (Latin 1) encoding is used
// for printable keys (such as A-Z, 0-9 etc), and values above 256
// represent Special (non-printable) keys (e.g. F1, Page Up etc).
const (
	KeyUnknown = -1
	KeySpace   = 32
	KeySpecial = 256
)

const (
	_ = (KeySpecial + iota)
	KeyEsc
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyLshift
	KeyRshift
	KeyLctrl
	KeyRctrl
	KeyLalt
	KeyRalt
	KeyTab
	KeyEnter
	KeyBackspace
	KeyInsert
	KeyDel
	KeyPageup
	KeyPagedown
	KeyHome
	KeyEnd
	KeyKP0
	KeyKP1
	KeyKP2
	KeyKP3
	KeyKP4
	KeyKP5
	KeyKP6
	KeyKP7
	KeyKP8
	KeyKP9
	KeyKPDidivde
	KeyKPMultiply
	KeyKPSubtract
	KeyKPAdd
	KeyKPDecimal
	KeyKPEqual
	KeyKPEnter
	KeyKPNumlock
	KeyCapslock
	KeyScrolllock
	KeyPause
	KeyLsuper
	KeyRsuper
	KeyMenu
	KeyLast = KeyMenu
)


