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
	glfw "github.com/go-gl/glfw3"
)

type EventHandler struct {
	MouseEvents chan *MouseEvent
}

func NewEventHandler(w *glfw.Window) (e *EventHandler) {
	e = &EventHandler{
		MouseEvents: make(chan *MouseEvent, 20),
	}
	w.SetCursorPositionCallback(e.onMouseMove)
	return
}

func (e *EventHandler) Poll() {
	glfw.PollEvents()
}

func (e *EventHandler) onMouseMove(w *glfw.Window, xoff float64, yoff float64) {
	event := &MouseEvent{
		X: float32(xoff),
		Y: float32(yoff),
	}
	select {
	case e.MouseEvents <- event:
		// Added to mouse events
	default:
		// Mouse events buffer is too full, not being read.
		// Drop the event on the floor.
		// TODO: Warn?
	}
}

type MouseEvent struct {
	X float32
	Y float32
}
