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

import ()

type Element struct {
	Children []Node
	parent   Node
	z        float64
	bounds   Rectangle
}

func (e *Element) AddChild(node Node) {
	node.SetParent(e)
	e.Children = append(e.Children, node)
}

func (e *Element) RemoveChild(node Node) {
	for i, c := range e.Children {
		if c == node {
			e.Children = append(e.Children[:i], e.Children[i+1:]...)
			break
		}
	}
	return
}

func (e *Element) GetAllChildren() []Node {
	r := append([]Node{}, e.Children[:]...)
	for _, c := range e.Children {
		r = append(r, c.GetAllChildren()[:]...)
	}
	return r
}

func (e *Element) Clear() {
	e.Children = make([]Node, 0)
}

func (e *Element) SetParent(node Node) {
	e.parent = node
}

func (e *Element) Parent() Node {
	return e.parent
}

func (e *Element) Draw() {
	for _, child := range e.Children {
		child.Draw()
	}
}

func (e *Element) SetBounds(r Rectangle) {
	e.bounds = r
}

func (e *Element) GlobalBounds() Rectangle {
	if e.parent == nil {
		return e.Bounds()
	}
	return e.Bounds().Add(e.parent.GlobalBounds().Min)
}

func (e *Element) Bounds() Rectangle {
	return e.bounds
}

func (e *Element) RelativeBounds(n Node) Rectangle {
	return e.GlobalBounds().Sub(n.GlobalBounds().Min)
}

func (e *Element) Move(p Point) {
	e.bounds = e.bounds.Add(p)
}

func (e *Element) MoveTo(p Point) {
	var (
		x = p.X + e.bounds.Max.X - e.bounds.Min.X
		y = p.Y + e.bounds.Max.Y - e.bounds.Min.Y
	)
	e.bounds = Rectangle{Min: p, Max: Pt(x, y)}
}

func (e *Element) Width() float64 {
	return e.bounds.Max.X - e.bounds.Min.X
}

func (e *Element) SetWidth(w float64) {
	e.bounds.Max.X = e.bounds.Min.X + w
}

func (e *Element) Height() float64 {
	return e.bounds.Max.Y - e.bounds.Min.Y
}

func (e *Element) SetHeight(h float64) {
	e.bounds.Max.Y = e.bounds.Min.Y + h
}

func (e *Element) X() float64 {
	return e.bounds.Min.X
}

func (e *Element) Y() float64 {
	return e.bounds.Min.Y
}

func (e *Element) Z() float64 {
	return e.z
}

func (e *Element) SetZ(z float64) {
	e.z = z
}
