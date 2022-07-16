// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package zone

import tea "github.com/charmbracelet/bubbletea"

// DefaultManager is a app-wide manager. To initialize it, call NewGlobal().
var DefaultManager *Manager

// NewGlobal initializes a global manager, so you don't have to pass the manager
// between all components. This is primarily only useful if you have full control
// of the zones you want to monitor, however if developing a library using this,
// make sure you allow users to pass in their own manager.
func NewGlobal() {
	if DefaultManager != nil {
		return
	}

	DefaultManager = New()
}

// Close stops the manager worker.
func Close() {
	DefaultManager.checkInitialized()
	DefaultManager.Close()
}

// Mark returns v wrapped with a start and end ANSI sequence to allow the zone
// manager to determine where the zone is, including its window offsets. The ANSI
// sequences used should be ignored by lipgloss width methods, to prevent incorrect
// width calculations.
func Mark(id, v string) string {
	DefaultManager.checkInitialized()
	return DefaultManager.Mark(id, v)
}

// Clear removes any stored zones for the given ID.
func Clear(id string) {
	DefaultManager.checkInitialized()
	DefaultManager.Clear(id)
}

// Get returns the zone info of the given ID. If the ID is not known (yet),
// Get() returns nil.
func Get(id string) (a *ZoneInfo) {
	DefaultManager.checkInitialized()
	return DefaultManager.Get(id)
}

// Scan will scan the view output, searching for zone markers, returning the
// original view output with the zone markers stripped. Scan() should be used
// by the outer most model/component of your application, and not inside of a
// model/component child.
//
// Scan buffers the zone info to be stored, so an immediate call to Get(id) may
// not return the correct information. Thus it's recommended to primarily use
// Get(id) for actions like mouse events, which don't occur immediately after a
// view shift (where the previously stored zone info might be different).
func Scan(v string) string {
	DefaultManager.checkInitialized()
	return DefaultManager.Scan(v)
}

// AnyInBounds sends a MsgZoneInBounds message to the provided model for each zone
// that is in the bounds of the provided mouse event. The results of the call to
// Update() are discarded.
//
// Note that if multiple zones are within bounds, each one will be sent as an event
// in alphabetical sorted order of the ID.
func AnyInBounds(model tea.Model, mouse tea.MouseMsg) {
	DefaultManager.checkInitialized()
	DefaultManager.AnyInBounds(model, mouse)
}

// NewPrefix generates a zone marker ID prefix, which can help prevent overlapping
// zone markers between multiple components. Each call to NewPrefix() returns a
// new unique prefix.
//
// Usage example:
//	func NewModel() tea.Model {
//		return &model{
//			id: zone.NewPrefix(),
//		}
//	}
//
//	type model struct {
//		id     string
//		active int
//		items  []string
//	}
//
//	func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//		switch msg := msg.(type) {
//		// [...]
//		case tea.MouseMsg:
//			// [...]
//			for i, item := range m.items {
//				if zone.Get(m.id + item.name).InBounds(msg) {
//					m.active = i
//					break
//				}
//			}
//		}
//		return m, nil
//	}
//
//	func (m model) View() string {
//		return zone.Mark(m.id+"some-other-id", "rendered stuff here")
//	}
func NewPrefix() string {
	DefaultManager.checkInitialized()
	return DefaultManager.NewPrefix()
}
