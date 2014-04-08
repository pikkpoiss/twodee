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
)

type MenuItemData struct {
	Key   int32
	Value int32
}

type MenuNode struct {
	data        *MenuItemData
	label       string
	parent      *MenuNode
	children    []*MenuNode
	highlighted bool
	back        bool
}

func (mn *MenuNode) IsHighlighted() bool {
	return mn.highlighted
}

func (mn *MenuNode) Label() string {
	return mn.label
}

type Menu struct {
	root        *MenuNode
	highlighted int
}

func NewMenuNode(key int32, val int32, label string, children []*MenuNode) *MenuNode {
	var node = &MenuNode{
		data:        &MenuItemData{Key: key, Value: val},
		label:       label,
		children:    children,
		back:        false,
		highlighted: false,
	}
	for _, child := range children {
		child.parent = node
	}
	return node
}

func BackMenuNode(label string) *MenuNode {
	return &MenuNode{
		data:        nil,
		label:       label,
		parent:      nil,
		children:    nil,
		highlighted: false,
		back:        true,
	}
}

func NewMenu(items []*MenuNode) (menu *Menu, err error) {
	if len(items) == 0 {
		err = fmt.Errorf("No items in menu")
		return
	}
	menu = &Menu{
		root: NewMenuNode(0, 0, "", items),
	}
	menu.updateHighlighted(0)
	return
}

func (m *Menu) getHighlighted() *MenuNode {
	return m.root.children[m.highlighted%len(m.root.children)]
}

func (m *Menu) Items() []*MenuNode {
	return m.root.children
}

func (m *Menu) Select() *MenuItemData {
	h := m.getHighlighted()
	if len(h.children) != 0 {
		m.root = h
		m.updateHighlighted(0)
		return nil
	} else if h.back == true && h.parent.parent != nil {
		m.root = h.parent.parent
		m.updateHighlighted(0)
		return nil
	}
	return h.data
}

func (m *Menu) Next() {
	m.updateHighlighted((m.highlighted + 1) % len(m.root.children))
}

func (m *Menu) Prev() {
	count := len(m.root.children)
	m.updateHighlighted((m.highlighted + count - 1) % count)
}

func (m *Menu) updateHighlighted(i int) {
	if len(m.root.children) <= i {
		return
	}
	m.root.children[m.highlighted].highlighted = false
	m.highlighted = i
	m.root.children[i].highlighted = true
}
