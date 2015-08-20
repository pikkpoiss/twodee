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
	"container/heap"
	"fmt"
)

type step struct {
	x        int32
	y        int32
	g        int32
	parent   *step
	priority int32 // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

type pathSet map[int32]*step

type pathQueue []*step

func (q pathQueue) Len() int { return len(q) }

func (q pathQueue) Less(i, j int) bool {
	return q[i].priority < q[j].priority
}

func (q pathQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *pathQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*step)
	item.index = n
	*q = append(*q, item)
}

func (q *pathQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*q = old[0 : n-1]
	return item
}

func (q pathQueue) Find(x, y int32) *step {
	for i := 0; i < len(q); i++ {
		if q[i].x == x && q[i].y == y {
			return q[i]
		}
	}
	return nil
}

func newStep(x, y, g, p int32, parent *step) *step {
	return &step{
		x:        x,
		y:        y,
		g:        g,
		parent:   parent,
		priority: p,
	}
}

// Pretty much a direct A* implementation from
// http://theory.stanford.edu/~amitp/GameProgramming/ImplementationNotes.html
func (g *Grid) GetPath(x1, y1, x2, y2 int32) (out []GridPoint, err error) {
	var (
		open     = &pathQueue{}
		closed   = pathSet{}
		current  *step
		xs       = []int32{0, 1, 0, -1}
		ys       = []int32{1, 0, -1, 0}
		cost     int32
		neighbor *step
	)
	// CLOSED = empty set
	// OPEN = priority queue containing START
	heap.Init(open)
	heap.Push(open, newStep(x1, y1, 0, 0, nil))

	// While lowest rank in OPEN is not the GOAL:
	for open.Len() > 0 {
		// Set current = remove lowest rank item from OPEN
		current = heap.Pop(open).(*step)
		if current.x == x2 && current.y == y2 {
			// Reconstruct reverse path from goal to start
			// by following parent pointers
			return getPoints(current), nil
		}
		// Add current to CLOSED
		closed[g.Index(current.x, current.y)] = current

		// For neighbors of current:
		for _, xm := range xs {
			for _, ym := range ys {
				// Set cost = g(current) + movementcost(current, neighbor)
				cost = current.g + 1
				var inopen, inclosed bool
				var nx, ny int32
				nx = current.x + xm
				ny = current.y + ym
				if g.Get(nx, ny) == true {
					// Spot is occupied
					continue
				}
				// If neighbor in OPEN and cost less than g(neighbor):
				// remove neighbor from OPEN, because new path is better
				if neighbor = open.Find(nx, ny); neighbor != nil {
					if cost < neighbor.g {
						heap.Remove(open, neighbor.index)
						inopen = false
					} else {
						inopen = true
					}
				} else {
					inopen = false
				}
				// If neighbor in CLOSED and cost less than g(neighbor):
				// remove neighbor from CLOSED
				if neighbor, inclosed = closed[g.Index(nx, ny)]; inclosed == true {
					if cost < neighbor.g {
						delete(closed, g.Index(nx, ny))
						inclosed = false
					}
				} else {
					inclosed = false
				}
				// If neighbor not in OPEN and neighbor not in CLOSED:
				if !inopen && !inclosed {
					var h = heuristic(nx, ny, x2, y2)
					// Set g(neighbor) to cost
					// Set priority queue rank to g(neighbor) + h(neighbor)
					// Set neighbor's parent to current
					neighbor = newStep(nx, ny, cost, cost+h, current)
					// Add neighbor to OPEN
					heap.Push(open, neighbor)
				}
			}
		}
	}
	err = fmt.Errorf("No path found")
	return
}

func heuristic(x1, y1, x2, y2 int32) int32 {
	var (
		dx   = x2 - x1
		dy   = y2 - y1
		negx = dx < 0
		negy = dy < 0
	)
	if negx == true {
		dx *= -1
	}
	if negy == true {
		dy *= -1
	}
	return dx + dy
}

type GridPoint struct {
	X, Y int32
}

func getPoints(dest *step) []GridPoint {
	var (
		count  = 1
		marker = dest
		out    []GridPoint
	)
	for marker.parent != nil {
		count += 1
		marker = marker.parent
	}
	out = make([]GridPoint, count)
	marker = dest
	for i := count - 1; i > -1; i-- {
		out[i] = GridPoint{marker.x, marker.y}
		marker = marker.parent
	}
	return out
}
