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

type MenuItem interface {
	Highlighted() bool
	Label() string
	setHighlighted(val bool)
	Parent() MenuItem
	setParent(item MenuItem)
	Active() bool
	setActive(val bool)
}

type baseMenuItem struct {
	label       string
	highlighted bool
	active      bool
	parent      MenuItem
}

func (mi *baseMenuItem) Label() string {
	return mi.label
}

func (mi *baseMenuItem) Highlighted() bool {
	return mi.highlighted
}

func (mi *baseMenuItem) setHighlighted(val bool) {
	mi.highlighted = val
}

func (mi *baseMenuItem) Active() bool {
	return mi.active
}

func (mi *baseMenuItem) setActive(val bool) {
	mi.active = val
}

func (mi *baseMenuItem) Parent() MenuItem {
	return mi.parent
}

func (mi *baseMenuItem) setParent(item MenuItem) {
	mi.parent = item
}

type BackMenuItem struct {
	baseMenuItem
}

func NewBackMenuItem(label string) (item *BackMenuItem) {
	item = &BackMenuItem{}
	item.baseMenuItem.label = label
	return
}

type ParentMenuItem struct {
	baseMenuItem
	children []MenuItem
}

func NewParentMenuItem(label string, children []MenuItem) (item *ParentMenuItem) {
	item = &ParentMenuItem{
		baseMenuItem: baseMenuItem{
			label: label,
		},
		children: children,
	}
	for _, child := range item.children {
		child.setParent(item)
	}
	return
}

func (mi *ParentMenuItem) Children() []MenuItem {
	return mi.children
}

type MenuItemData struct {
	Key   int32
	Value int32
}

type KeyValueMenuItem struct {
	baseMenuItem
	data *MenuItemData
}

func NewKeyValueMenuItem(label string, key, value int32) *KeyValueMenuItem {
	return &KeyValueMenuItem{
		baseMenuItem: baseMenuItem{label: label},
		data: &MenuItemData{
			Key:   key,
			Value: value,
		},
	}
}

func (mi *KeyValueMenuItem) Data() *MenuItemData {
	return mi.data
}

type BoundValueMenuItem struct {
	KeyValueMenuItem
	dest *int32
}

func NewBoundValueMenuItem(label string, val int32, dest *int32) *BoundValueMenuItem {
	return &BoundValueMenuItem{
		KeyValueMenuItem: *NewKeyValueMenuItem(label, -1, val),
		dest:             dest,
	}
}

func (mi *BoundValueMenuItem) IsSameAsDest() bool {
	return *(mi.dest) == mi.KeyValueMenuItem.data.Value
}

func (mi *BoundValueMenuItem) SetDest() {
	*(mi.dest) = mi.KeyValueMenuItem.data.Value
}

type Menu struct {
	root        *ParentMenuItem
	items       []MenuItem
	highlighted int
}

func NewMenu(items []MenuItem) (menu *Menu, err error) {
	if len(items) == 0 {
		err = fmt.Errorf("No items in menu")
		return
	}
	menu = &Menu{
		root:  NewParentMenuItem("root", items),
		items: items,
	}
	menu.updateHighlighted(0)
	return
}

func (m *Menu) getHighlighted() MenuItem {
	return m.items[m.highlighted%len(m.items)]
}

func (m *Menu) Items() []MenuItem {
	return m.items
}

func (m *Menu) Reset() {
	m.items = m.root.Children()
	m.updateHighlighted(0)
}

func (m *Menu) SelectItem(item MenuItem) *MenuItemData {
	if item == nil {
		return nil
	}
	switch selected := item.(type) {
	case *BackMenuItem:
		m.SelectItem(selected.Parent().Parent())
	case *ParentMenuItem:
		m.items = selected.Children()
		m.updateHighlighted(0)
	case *KeyValueMenuItem:
		return selected.Data()
	case *BoundValueMenuItem:
		selected.SetDest()
		m.updateHighlighted(m.highlighted)
	}
	return nil
}

func (m *Menu) HighlightItem(h MenuItem) {
	for i, item := range m.items {
		if item == h {
			m.updateHighlighted(i)
			break
		}
	}
}

func (m *Menu) Select() *MenuItemData {
	return m.SelectItem(m.getHighlighted())
}

func (m *Menu) Next() {
	m.updateHighlighted((m.highlighted + 1) % len(m.items))
}

func (m *Menu) Prev() {
	count := len(m.items)
	m.updateHighlighted((m.highlighted + count - 1) % count)
}

func (m *Menu) updateHighlighted(i int) {
	for _, item := range m.items {
		switch curr := item.(type) {
		case *BoundValueMenuItem:
			item.setActive(curr.IsSameAsDest())
		}
	}
	if len(m.items) <= i {
		return
	}
	if len(m.items) > m.highlighted {
		m.items[m.highlighted].setHighlighted(false)
	}
	m.highlighted = i
	m.items[i].setHighlighted(true)
}
