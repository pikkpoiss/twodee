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
	"github.com/go-gl/glfw/v3.1/glfw"
)

type EventHandler struct {
	Events chan Event
}

func NewEventHandler(w *glfw.Window) (e *EventHandler) {
	e = &EventHandler{
		Events: make(chan Event, 100),
	}
	w.SetCursorPositionCallback(e.onMouseMove)
	w.SetKeyCallback(e.onKey)
	w.SetMouseButtonCallback(e.onMouseButton)
	return
}

func (e *EventHandler) Poll() {
	glfw.PollEvents()
}

func (e *EventHandler) onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	event := &KeyEvent{
		Code: KeyCode(key),
		Type: Action(action),
	}
	e.enqueue(event)
}

func (e *EventHandler) onMouseMove(w *glfw.Window, xoff float64, yoff float64) {
	event := &MouseMoveEvent{
		X: float32(xoff),
		Y: float32(yoff),
	}
	e.enqueue(event)
}

func (e *EventHandler) onMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	event := &MouseButtonEvent{
		Type:   Action(action),
		Button: MouseButton(button),
	}
	e.enqueue(event)
}

func (e *EventHandler) enqueue(event Event) {
	select {
	case e.Events <- event:
		// Added to events
	default:
		// Events buffer is too full, not being read.
		// Drop the event on the floor.
		// TODO: Warn?
	}
}

type Event interface{}

type MouseMoveEvent struct {
	X float32
	Y float32
}

type MouseButtonEvent struct {
	Button MouseButton
	Type   Action
}

type KeyEvent struct {
	Code KeyCode
	Type Action
}

type KeyCode int

const (
	KeyUp     = KeyCode(glfw.KeyUp)
	KeyDown   = KeyCode(glfw.KeyDown)
	KeyLeft   = KeyCode(glfw.KeyLeft)
	KeyRight  = KeyCode(glfw.KeyRight)
	KeyEnter  = KeyCode(glfw.KeyEnter)
	KeyEscape = KeyCode(glfw.KeyEscape)
	KeySpace  = KeyCode(glfw.KeySpace)
	KeyA      = KeyCode(glfw.KeyA)
	KeyB      = KeyCode(glfw.KeyB)
	KeyC      = KeyCode(glfw.KeyC)
	KeyD      = KeyCode(glfw.KeyD)
	KeyE      = KeyCode(glfw.KeyE)
	KeyF      = KeyCode(glfw.KeyF)
	KeyG      = KeyCode(glfw.KeyG)
	KeyH      = KeyCode(glfw.KeyH)
	KeyI      = KeyCode(glfw.KeyI)
	KeyJ      = KeyCode(glfw.KeyJ)
	KeyK      = KeyCode(glfw.KeyK)
	KeyL      = KeyCode(glfw.KeyL)
	KeyM      = KeyCode(glfw.KeyM)
	KeyN      = KeyCode(glfw.KeyN)
	KeyO      = KeyCode(glfw.KeyO)
	KeyP      = KeyCode(glfw.KeyP)
	KeyQ      = KeyCode(glfw.KeyQ)
	KeyR      = KeyCode(glfw.KeyR)
	KeyS      = KeyCode(glfw.KeyS)
	KeyT      = KeyCode(glfw.KeyT)
	KeyU      = KeyCode(glfw.KeyU)
	KeyV      = KeyCode(glfw.KeyV)
	KeyW      = KeyCode(glfw.KeyW)
	KeyX      = KeyCode(glfw.KeyX)
	KeyY      = KeyCode(glfw.KeyY)
	KeyZ      = KeyCode(glfw.KeyZ)
)

type Action int

const (
	Release = Action(glfw.Release)
	Press   = Action(glfw.Press)
	Repeat  = Action(glfw.Repeat)
)

type MouseButton int

const (
	MouseButton1      = MouseButton(glfw.MouseButton1)
	MouseButton2      = MouseButton(glfw.MouseButton2)
	MouseButton3      = MouseButton(glfw.MouseButton3)
	MouseButton4      = MouseButton(glfw.MouseButton4)
	MouseButton5      = MouseButton(glfw.MouseButton5)
	MouseButton6      = MouseButton(glfw.MouseButton6)
	MouseButton7      = MouseButton(glfw.MouseButton7)
	MouseButton8      = MouseButton(glfw.MouseButton8)
	MouseButtonLast   = MouseButton(glfw.MouseButtonLast)
	MouseButtonLeft   = MouseButton(glfw.MouseButtonLeft)
	MouseButtonRight  = MouseButton(glfw.MouseButtonRight)
	MouseButtonMiddle = MouseButton(glfw.MouseButtonMiddle)
)
