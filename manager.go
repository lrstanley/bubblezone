// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package zone

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// Have to use ansi escape codes to ensure lipgloss doesn't consider ID's as
	// part of the width of the view.
	identStart   = '\x1B' // ANSI escape code.
	identBracket = '['
	identEnd     = 'Z' // escape terminator.
)

var (
	markerCounter int64 = 1000 // Protected by atomic operations.
	prefixCounter int64 = 0    // Protected by atomic operations.
)

func New() (m *Manager) {
	m = &Manager{
		setChan: make(chan *ZoneInfo, 200),
		zones:   make(map[string]*ZoneInfo),
		ids:     make(map[string]string),
		rids:    make(map[string]string),
	}

	m.ctx, m.cancel = context.WithCancel(context.Background())
	go m.zoneWorker()

	return m
}

// Manager holds the state of the zone manager, including ID zones and
// zones of components.
type Manager struct {
	ctx     context.Context
	cancel  func()
	setChan chan *ZoneInfo

	zoneMu sync.RWMutex
	zones  map[string]*ZoneInfo

	idMu sync.RWMutex
	ids  map[string]string // user ID -> generated control sequence ID.
	rids map[string]string // generated control sequence ID -> user ID.
}

func (m *Manager) checkInitialized() {
	if m == nil {
		panic("manager not initialized")
	}
}

// Close stops the manager worker.
func (m *Manager) Close() {
	m.cancel()
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
func (m *Manager) NewPrefix() string {
	return "zone_" + strconv.FormatInt(atomic.AddInt64(&prefixCounter, 1), 10) + "__"
}

// Mark returns v wrapped with a start and end ANSI sequence to allow the zone
// manager to determine where the zone is, including its window offsets. The ANSI
// sequences used should be ignored by lipgloss width methods, to prevent incorrect
// width calculations.
func (m *Manager) Mark(id, v string) string {
	if id == "" || v == "" {
		return v
	}

	m.idMu.RLock()
	gid := m.ids[id]
	m.idMu.RUnlock()

	if gid != "" {
		return gid + v + gid
	}

	m.idMu.Lock()
	gid = string(identStart) + string(identBracket) + strconv.FormatInt(atomic.AddInt64(&markerCounter, 1), 10) + string(identEnd)
	m.ids[id] = gid
	m.rids[gid] = id
	m.idMu.Unlock()

	return gid + v + gid
}

// Clear removes any stored zones for the given ID.
func (m *Manager) Clear(id string) {
	m.zoneMu.Lock()
	delete(m.zones, id)
	m.zoneMu.Unlock()
}

// Get returns the zone info of the given ID. If the ID is not known (yet),
// Get() returns nil.
func (m *Manager) Get(id string) (zone *ZoneInfo) {
	m.zoneMu.RLock()
	zone = m.zones[id]
	m.zoneMu.RUnlock()
	return zone
}

// getReverse returns the component ID from a generated ID (that includes ANSI
// escape codes).
func (m *Manager) getReverse(id string) (resolved string) {
	m.idMu.RLock()
	resolved = m.rids[id]
	m.idMu.RUnlock()
	return resolved
}

func (m *Manager) zoneWorker() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case xy := <-m.setChan:
			m.zoneMu.Lock()
			if xy.id != "" {
				m.zones[m.getReverse(xy.id)] = xy
			} else {
				// Assume previous iterations are cleared.
				for k := range m.zones {
					if m.zones[k].iteration != xy.iteration {
						delete(m.zones, k)
					}
				}
			}
			m.zoneMu.Unlock()
		}
	}
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
func (m *Manager) Scan(v string) string {
	iteration := time.Now().Nanosecond()
	s := newScanner(m, v, iteration)
	s.run()
	m.setChan <- &ZoneInfo{iteration: iteration}
	return s.input
}
