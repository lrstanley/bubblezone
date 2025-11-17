// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package zone

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

var _ tea.Model = (*testModel)(nil)

type testModel struct {
	received []tea.Msg
}

func newTestModel() *testModel {
	return &testModel{}
}

func (m *testModel) Init() tea.Cmd {
	return nil
}

func (m *testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		go AnyInBounds(m, msg)
		return m, nil
	case MsgZoneInBounds:
		m.received = append(m.received, msg)
	}
	return m, nil
}

func (m *testModel) View() tea.View {
	return tea.NewView(Scan("test\nfoo\naaa " + Mark("foo", "bar\ntest123456789") + " aaa\nbaz"))
}

func TestAnyInBounds(t *testing.T) {
	m := newTestModel()
	_ = m.View()
	time.Sleep(100 * time.Millisecond)
	xy := Get("foo")
	if xy.IsZero() {
		t.Error("id not found")
	}

	_, _ = m.Update(tea.MouseMotionMsg{X: 4, Y: 2})
	time.Sleep(100 * time.Millisecond)

	var contains bool
	for _, msg := range m.received {
		if evt, ok := msg.(MsgZoneInBounds); ok {
			if evt.Zone.id == xy.id {
				contains = true
				break
			}
		}
	}

	if !contains {
		t.Error("expected true")
	}
}

var _ tea.Model = (*testModelValue)(nil)

type testModelValue struct {
	received []tea.Msg
}

func (m testModelValue) Init() tea.Cmd {
	return nil
}

func (m testModelValue) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		return AnyInBoundsAndUpdate(m, msg)
	case MsgZoneInBounds:
		m.received = append(m.received, msg)
	}
	return m, nil
}

func (m testModelValue) View() tea.View {
	return tea.NewView(Scan("test\nfoo\naaa " + Mark("foo", "bar\ntest123456789") + " aaa\nbaz"))
}

func TestAnyInBoundsAndUpdate(t *testing.T) {
	m := testModelValue{}

	_ = m.View()
	time.Sleep(100 * time.Millisecond)
	xy := Get("foo")
	if xy.IsZero() {
		t.Error("id not found")
	}

	newModel, _ := m.Update(tea.MouseMotionMsg{X: 4, Y: 2})
	m = newModel.(testModelValue)
	time.Sleep(100 * time.Millisecond)

	var contains bool

	for _, msg := range m.received {
		if evt, ok := msg.(MsgZoneInBounds); ok {
			if evt.Zone.id == xy.id {
				contains = true
				break
			}
		}
	}

	if !contains {
		t.Error("expected true")
	}
}
