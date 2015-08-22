// Copyright 2015 Arne Roomann-Kurrik
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
	"testing"
	"time"
)

// Tests that setting a callback works and the callback is called.
func TestBoundedAnimationSetCallback(t *testing.T) {
	var (
		done = false
		anim = BoundedAnimation{0, 1 * time.Second, nil}
		cb   = func() { done = true }
	)
	anim.SetCallback(cb)
	anim.Update(1 * time.Second)
	if done != true {
		t.Fatalf("BoundedAnimation callback not called")
	}
}

// Tests that BoundedAnimation.IsDone is true when elapsed == duration.
func TestBoundedAnimationIsDoneElapsedEqDuration(t *testing.T) {
	var anim = BoundedAnimation{0, 1 * time.Second, nil}
	anim.Update(1 * time.Second)
	if !anim.IsDone() {
		t.Fatalf("BoundedAnimation.IsDone was not expected value")
	}
}

// Tests that BoundedAnimation.IsDone is true when elapsed > duration.
func TestBoundedAnimationIsDoneElapsedGtDuration(t *testing.T) {
	var anim = BoundedAnimation{0, 1 * time.Second, nil}
	anim.Update(9 * time.Second)
	if !anim.IsDone() {
		t.Fatalf("BoundedAnimation.IsDone was not expected value")
	}
}

// Tests that BoundedAnimation.IsDone is false when elapsed < duration.
func TestBoundedAnimationIsDoneElapsedLtDuration(t *testing.T) {
	var anim = BoundedAnimation{0, 1 * time.Second, nil}
	anim.Update(200 * time.Millisecond)
	if anim.IsDone() {
		t.Fatalf("BoundedAnimation.IsDone was not expected value")
	}
}

// Tests that BoundedAnimation.Update increments elapsed and returns overflow.
func TestBoundedAnimationUpdate(t *testing.T) {
	var (
		anim               = BoundedAnimation{0, 1 * time.Second, nil}
		resp time.Duration = 0
	)
	resp = anim.Update(200 * time.Millisecond)
	if anim.Elapsed != 200*time.Millisecond {
		t.Fatalf("BoundedAnimation.Elapsed was not expected value")
	}
	if resp != 0 {
		t.Fatalf("BoundedAnimation.Elapsed did not return expected value")
	}
	resp = anim.Update(600 * time.Millisecond)
	if anim.Elapsed != 800*time.Millisecond {
		t.Fatalf("BoundedAnimation.Elapsed was not expected value")
	}
	if resp != 0 {
		t.Fatalf("BoundedAnimation.Elapsed did not return expected value")
	}
	resp = anim.Update(400 * time.Millisecond)
	if anim.Elapsed != 1200*time.Millisecond {
		t.Fatalf("BoundedAnimation.Elapsed was not expected value")
	}
	if resp != 200*time.Millisecond {
		t.Fatalf("BoundedAnimation.Elapsed did not return expected value")
	}
}

func TestBoundedAnimationReset(t *testing.T) {
	var anim = BoundedAnimation{200 * time.Millisecond, 1 * time.Second, nil}
	anim.Reset()
	if anim.Elapsed != 0 {
		t.Fatalf("BoundedAnimation.Elapsed was not expected value")
	}
}

func TestGroupedAnimationSetCallback(t *testing.T) {
	var (
		done   = false
		child1 = &BoundedAnimation{0, 1 * time.Second, nil}
		child2 = &BoundedAnimation{0, 2 * time.Second, nil}
		anim   = GroupedAnimation{[]Animator{child1, child2}, nil}
		cb     = func() { done = true }
	)
	anim.SetCallback(cb)
	anim.Update(1 * time.Second)
	if done {
		t.Fatalf("GroupedAnimation callback called too early")
	}
	if !child1.IsDone() {
		t.Fatalf("First child animation not marked done")
	}
	if child2.IsDone() {
		t.Fatalf("Second child animation marked done too early")
	}
	anim.Update(1 * time.Second)
	if !done {
		t.Fatalf("GroupedAnimation callback not called")
	}
	if !child2.IsDone() {
		t.Fatalf("Second child animation not marked done")
	}
	if child1.Elapsed != 2*time.Second {
		t.Fatalf("First child elapsed not equal to expected value")
	}
	if child2.Elapsed != 2*time.Second {
		t.Fatalf("Second child elapsed not equal to expected value")
	}
}

func TestGroupedAnimationIsDone(t *testing.T) {
	var (
		child1 = &BoundedAnimation{0, 1 * time.Second, nil}
		child2 = &BoundedAnimation{0, 2 * time.Second, nil}
		anim   = GroupedAnimation{[]Animator{child1, child2}, nil}
	)
	anim.Update(1 * time.Second)
	if anim.IsDone() {
		t.Fatalf("GroupedAnimation.IsDone is true too early")
	}
	anim.Update(1 * time.Second)
	if !anim.IsDone() {
		t.Fatalf("GroupedAnimation.IsDone is not true after animation")
	}
}

func TestGroupedAnimationUpdate(t *testing.T) {
	var (
		child1               = &BoundedAnimation{0, 1 * time.Second, nil}
		child2               = &BoundedAnimation{0, 2 * time.Second, nil}
		anim                 = GroupedAnimation{[]Animator{child1, child2}, nil}
		resp   time.Duration = 0
	)
	resp = anim.Update(1 * time.Second)
	if resp != 0 {
		t.Fatalf("GroupedAnimation.Update did not return correct remainder")
	}
	resp = anim.Update(200 * time.Millisecond)
	if resp != 0 {
		t.Fatalf("GroupedAnimation.Update did not return correct remainder")
	}
	resp = anim.Update(900 * time.Millisecond)
	if resp != 100*time.Millisecond {
		t.Fatalf("GroupedAnimation.Update did not return correct remainder, got %v", resp)
	}
}

func TestGroupedAnimationReset(t *testing.T) {
	var (
		child1 = &BoundedAnimation{200 * time.Millisecond, 1 * time.Second, nil}
		child2 = &BoundedAnimation{1200 * time.Millisecond, 2 * time.Second, nil}
		anim   = GroupedAnimation{[]Animator{child1, child2}, nil}
	)
	anim.Reset()
	if child1.Elapsed != 0 {
		t.Fatalf("GroupedAnimation.Reset did not reset first child")
	}
	if child2.Elapsed != 0 {
		t.Fatalf("GroupedAnimation.Reset did not reset second child")
	}
}

func TestGroupedAnimationDelete(t *testing.T) {
	var (
		child1 = &BoundedAnimation{0, 1 * time.Second, nil}
		child2 = &BoundedAnimation{0, 2 * time.Second, nil}
		anim   = GroupedAnimation{[]Animator{child1, child2}, nil}
	)
	anim.Delete()
	anim.Update(100 * time.Millisecond)
	if child1.Elapsed != 0 {
		t.Fatalf("GroupedAnimation.Delete must remove references to children")
	}
}

func TestChainedAnimationSetCallback(t *testing.T) {
	var (
		done   = false
		child1 = &BoundedAnimation{0, 1 * time.Second, nil}
		child2 = &BoundedAnimation{0, 2 * time.Second, nil}
		anim   = ChainedAnimation{[]Animator{child1, child2}, false, 0, nil}
		cb     = func() { done = true }
	)
	anim.SetCallback(cb)
	anim.Update(1 * time.Second)
	if done {
		t.Fatalf("ChainedAnimation callback called too early")
	}
	anim.Update(2 * time.Second)
	if !done {
		t.Fatalf("ChainedAnimation callback not called")
	}
}

func TestChainedAnimationIsDoneNoLoop(t *testing.T) {
	var (
		child1 = &BoundedAnimation{0, 1 * time.Second, nil}
		child2 = &BoundedAnimation{0, 2 * time.Second, nil}
		anim   = ChainedAnimation{[]Animator{child1, child2}, false, 0, nil}
	)
	if anim.IsDone() {
		t.Fatalf("ChainedAnimation.IsDone true too early")
	}
	anim.Update(1 * time.Second)
	if anim.IsDone() {
		t.Fatalf("ChainedAnimation.IsDone true too early")
	}
	anim.Update(1 * time.Second)
	if anim.IsDone() {
		t.Fatalf("ChainedAnimation.IsDone true too early")
	}
	anim.Update(1 * time.Second)
	if !anim.IsDone() {
		t.Fatalf("ChainedAnimation.IsDone not set after animation finishes")
	}
}

func TestChainedAnimationIsDoneLoop(t *testing.T) {
	var (
		child1 = &BoundedAnimation{0, 1 * time.Second, nil}
		child2 = &BoundedAnimation{0, 2 * time.Second, nil}
		anim   = ChainedAnimation{[]Animator{child1, child2}, true, 0, nil}
	)
	if anim.IsDone() {
		t.Fatalf("ChainedAnimation.IsDone should not be true for loops")
	}
	anim.Update(1 * time.Second)
	if anim.IsDone() {
		t.Fatalf("ChainedAnimation.IsDone should not be true for loops")
	}
	anim.Update(1 * time.Second)
	if anim.IsDone() {
		t.Fatalf("ChainedAnimation.IsDone should not be true for loops")
	}
	anim.Update(1 * time.Second)
	if anim.IsDone() {
		t.Fatalf("ChainedAnimation.IsDone should not be true for loops")
	}
	anim.Update(1 * time.Second)
}

func TestChainedAnimationUpdateNoLoop(t *testing.T) {
	var (
		child1 = &BoundedAnimation{0, 1 * time.Second, nil}
		child2 = &BoundedAnimation{0, 2 * time.Second, nil}
		anim   = ChainedAnimation{[]Animator{child1, child2}, false, 0, nil}
		resp   time.Duration
	)
	resp = anim.Update(500 * time.Millisecond)
	if child1.Elapsed != 500*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must update current child")
	}
	if resp != 0 {
		t.Fatalf("ChainedAnimation.Update must return 0 before animation is done")
	}
	resp = anim.Update(600 * time.Millisecond)
	if child1.Elapsed != 1100*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must update current child")
	}
	if child2.Elapsed != 100*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must move to second child if first overflows")
	}
	if resp != 0 {
		t.Fatalf("ChainedAnimation.Update must return 0 before animation is done")
	}
	resp = anim.Update(500 * time.Millisecond)
	if child1.Elapsed != 1100*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must not update previous child")
	}
	if child2.Elapsed != 600*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must update current child")
	}
	if resp != 0 {
		t.Fatalf("ChainedAnimation.Update must return 0 before animation is done")
	}
	resp = anim.Update(1500 * time.Millisecond)
	if child1.Elapsed != 1100*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must update current child")
	}
	if child2.Elapsed != 2100*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must update current child")
	}
	if resp != 100*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must return remainder when animation is done")
	}
}

func TestChainedAnimationUpdateLoop(t *testing.T) {
	var (
		child1 = &BoundedAnimation{0, 1 * time.Second, nil}
		child2 = &BoundedAnimation{0, 2 * time.Second, nil}
		anim   = ChainedAnimation{[]Animator{child1, child2}, true, 0, nil}
		resp   time.Duration
	)
	resp = anim.Update(500 * time.Millisecond)
	if child1.Elapsed != 500*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must update current child")
	}
	if resp != 0 {
		t.Fatalf("ChainedAnimation.Update must return 0 for loops")
	}
	resp = anim.Update(600 * time.Millisecond)
	if child1.Elapsed != 0 {
		t.Fatalf("ChainedAnimation.Update must reset previous child")
	}
	if child2.Elapsed != 100*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must move to second child if first overflows")
	}
	if resp != 0 {
		t.Fatalf("ChainedAnimation.Update must return 0 for loops")
	}
	resp = anim.Update(500 * time.Millisecond)
	if child1.Elapsed != 0 {
		t.Fatalf("ChainedAnimation.Update must not update previous child")
	}
	if child2.Elapsed != 600*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must update current child")
	}
	if resp != 0 {
		t.Fatalf("ChainedAnimation.Update must return 0 for loops")
	}
	resp = anim.Update(1400 * time.Millisecond)
	if child1.Elapsed != 0 {
		t.Fatalf("ChainedAnimation.Update must not update previous child")
	}
	if child2.Elapsed != 0 {
		t.Fatalf("ChainedAnimation.Update must reset child when done")
	}
	if resp != 0 {
		t.Fatalf("ChainedAnimation.Update must return 0 for loops")
	}
	resp = anim.Update(300 * time.Millisecond)
	if child1.Elapsed != 300*time.Millisecond {
		t.Fatalf("ChainedAnimation.Update must loop back to first child")
	}
	if child2.Elapsed != 0 {
		t.Fatalf("ChainedAnimation.Update must update current child")
	}
	if resp != 0 {
		t.Fatalf("ChainedAnimation.Update must return 0 for loops")
	}
}

func TestChainedAnimationReset(t *testing.T) {
	var (
		child1 = &BoundedAnimation{100 * time.Millisecond, 1 * time.Second, nil}
		child2 = &BoundedAnimation{200 * time.Millisecond, 2 * time.Second, nil}
		anim   = ChainedAnimation{[]Animator{child1, child2}, true, 0, nil}
	)
	anim.Reset()
	if child1.Elapsed != 0 {
		t.Fatalf("ChainedAnimation.Reset must reset first child")
	}
	if child2.Elapsed != 0 {
		t.Fatalf("ChainedAnimation.Reset must reset second child")
	}
}

func TestChainedAnimationDelete(t *testing.T) {
	var (
		child1 = &BoundedAnimation{0, 1 * time.Second, nil}
		child2 = &BoundedAnimation{0, 2 * time.Second, nil}
		anim   = ChainedAnimation{[]Animator{child1, child2}, true, 0, nil}
	)
	anim.Delete()
	anim.Update(100 * time.Millisecond)
	if child1.Elapsed != 0 {
		t.Fatalf("ChainedAnimation.Delete must remove references to children")
	}
}

func TestLinearAnimation(t *testing.T) {
	var (
		dest float32 = 1.0
		anim         = NewLinearAnimation(&dest, 10, 20, 5*time.Second)
	)
	if dest != 1.0 {
		t.Fatalf("Target value does not match expected")
	}
	anim.Update(1000 * time.Millisecond)
	if dest != 12 {
		t.Fatalf("Target value does not match expected, got %v", dest)
	}
	anim.Update(1000 * time.Millisecond)
	if dest != 14 {
		t.Fatalf("Target value does not match expected, got %v", dest)
	}
	anim.Update(1000 * time.Millisecond)
	if dest != 16 {
		t.Fatalf("Target value does not match expected, got %v", dest)
	}
	anim.Update(1000 * time.Millisecond)
	if dest != 18 {
		t.Fatalf("Target value does not match expected, got %v", dest)
	}
	anim.Update(1000 * time.Millisecond)
	if dest != 20 {
		t.Fatalf("Target value does not match expected, got %v", dest)
	}
	anim.Update(1000 * time.Millisecond)
	if dest != 20 {
		t.Fatalf("Target value does not match expected, got %v", dest)
	}
}
