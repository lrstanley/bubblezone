// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package zone

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestValidPosition(t *testing.T) {
	// Starts at X:4, Y:2, ends at X:12, Y:3.
	_ = Scan("test\nfoo\naaa " + Mark("foo", "bar\ntest123456789") + " aaa\nbaz")
	time.Sleep(100 * time.Millisecond)
	xy := Get("foo")
	if xy.IsZero() {
		t.Error("id not found")
	}

	if xy.StartX != 4 || xy.StartY != 2 || xy.EndX != 12 || xy.EndY != 3 {
		t.Errorf("got %#v, want %#v", xy, &ZoneInfo{
			id:        xy.id,
			iteration: xy.iteration,
			StartX:    4,
			StartY:    2,
			EndX:      12,
			EndY:      3,
		})
	}
}

func TestInBounds(t *testing.T) {
	// Starts at X:4, Y:2, ends at X:12, Y:3.
	_ = Scan("test\nfoo\naaa " + Mark("foo", "bar\ntest123456789") + " aaa\nbaz")
	time.Sleep(100 * time.Millisecond)
	xy := Get("foo")
	if xy.IsZero() {
		t.Error("id not found")
	}

	// Outside left.
	if xy.InBounds(tea.MouseMsg{X: 0, Y: 0}) {
		t.Error("expected false")
	}

	// Outside directly left.
	if xy.InBounds(tea.MouseMsg{X: 3, Y: 3}) {
		t.Error("expected false")
	}

	// Outside right.
	if xy.InBounds(tea.MouseMsg{X: 99, Y: 99}) {
		t.Error("expected false")
	}

	// Outside directly right.
	if xy.InBounds(tea.MouseMsg{X: 13, Y: 3}) {
		t.Error("expected false")
	}

	// Outside top.
	if xy.InBounds(tea.MouseMsg{X: 4, Y: 1}) {
		t.Error("expected false")
	}

	// Outside bottom.
	if xy.InBounds(tea.MouseMsg{X: 4, Y: 4}) {
		t.Error("expected false")
	}

	// Inside left top.
	if !xy.InBounds(tea.MouseMsg{X: 4, Y: 2}) {
		t.Error("expected true")
	}

	// Inside right bottom.
	if !xy.InBounds(tea.MouseMsg{X: 12, Y: 3}) {
		t.Error("expected true")
	}

	_ = Scan("test " + Mark("foo", "bar\nt") + " other things here")
	time.Sleep(100 * time.Millisecond)

	xy = Get("foo")
	if xy.InBounds(tea.MouseMsg{X: 2, Y: 1}) {
		t.Error("expected false")
	}
}

func TestInBoundsZero(t *testing.T) {
	xy := &ZoneInfo{}
	if xy.InBounds(tea.MouseMsg{X: 0, Y: 0}) {
		t.Error("expected false")
	}

	xy = Get("non-existent")
	if xy.InBounds(tea.MouseMsg{X: 0, Y: 0}) {
		t.Error("expected false")
	}
}

func TestPos(t *testing.T) {
	// Starts at X:4, Y:2, ends at X:12, Y:3.
	_ = Scan("test\nfoo\naaa " + Mark("foo", "bar\ntest123456789") + " aaa\nbaz")
	time.Sleep(100 * time.Millisecond)
	xy := Get("foo")
	if xy.IsZero() {
		t.Error("id not found")
	}

	if x, y := xy.Pos(tea.MouseMsg{X: 4, Y: 2}); x != 0 || y != 0 {
		t.Error("expected 0, 0")
	}

	if x, y := xy.Pos(tea.MouseMsg{X: 5, Y: 2}); x != 1 || y != 0 {
		t.Error("expected 1, 0")
	}

	xy = &ZoneInfo{}
	if x, y := xy.Pos(tea.MouseMsg{X: 0, Y: 0}); x != -1 || y != -1 {
		t.Error("expected -1, -1")
	}

	xy = Get("non-existent")
	if x, y := xy.Pos(tea.MouseMsg{X: 0, Y: 0}); x != -1 || y != -1 {
		t.Error("expected -1, -1")
	}
}
