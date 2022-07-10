// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package zone

import (
	"context"
	"strconv"
	"sync"
	"unicode/utf8"

	"github.com/muesli/ansi"
)

const (
	// Have to use ansi escape codes to ensure lipgloss doesn't consider ID's as
	// part of the width of the view.
	identStart    = '\x1B' // ANSI escape code.
	identStartLen = len(string(identStart))
	identEnd      = '\x9C' // ANSI termination code.
	identEndLen   = len(string(identEnd))

	areaStart = "__start"
	areaEnd   = "__end"
)

func New(beginCounterAt int) (m *Manager) {
	if beginCounterAt == 0 {
		beginCounterAt = 500
	}

	if beginCounterAt < 500 {
		panic("beginCounterAt must be >= 500 to prevent collisions with standard ANSI sequences")
	}

	m = &Manager{
		setChan:   make(chan *Position, 200),
		mapping:   make(map[string]*Position),
		ids:       make(map[string]string),
		rids:      make(map[string]string),
		idCounter: beginCounterAt,
	}

	m.ctx, m.cancel = context.WithCancel(context.Background())
	go m.worker()

	return m
}

// Manager holds the state of the zone manager, including ID mappings and
// zones of components.
type Manager struct {
	ctx     context.Context
	cancel  func()
	setChan chan *Position

	mapMu   sync.RWMutex
	mapping map[string]*Position

	idMu      sync.RWMutex
	idCounter int
	ids       map[string]string // user ID -> generated control sequence ID.
	rids      map[string]string // generated control sequence ID -> user ID.
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

// Mark returns v wrapped with a start and end ANSI sequence to allow the zone
// manager to determine where the zone is, including its window offsets. The ANSI
// sequences used should be ignored by lipgloss width methods, to prevent incorrect
// width calculations.
func (m *Manager) Mark(id, v string) string {
	startID := id + areaStart
	endID := id + areaEnd

	m.idMu.RLock()
	start := m.ids[startID]
	end := m.ids[endID]
	m.idMu.RUnlock()

	if start != "" && end != "" {
		return start + v + end
	}

	m.idMu.Lock()

	m.idCounter++
	counter := strconv.Itoa(m.idCounter)
	start = string(identStart) + counter + string(identEnd)
	m.ids[startID] = start
	m.rids[counter] = startID // TODO: should this be counter, or start?

	m.idCounter++
	counter = strconv.Itoa(m.idCounter)
	end = string(identStart) + counter + string(identEnd)
	m.ids[endID] = end
	m.rids[counter] = endID

	m.idMu.Unlock()
	return start + v + end
}

// Clear removes any stored zones for the given ID.
func (m *Manager) Clear(id string) {
	m.mapMu.Lock()
	delete(m.mapping, id+areaStart)
	delete(m.mapping, id+areaEnd)
	m.mapMu.Unlock()
}

// Get returns the zone info of the given ID. If the ID is not known (yet),
// Get() returns nil.
func (m *Manager) Get(id string) (a *ZoneInfo) {
	m.mapMu.RLock()
	a = &ZoneInfo{
		Start: m.mapping[id+areaStart],
		End:   m.mapping[id+areaEnd],
	}
	m.mapMu.RUnlock()
	return a
}

// getReverse returns the component ID from a generated ID (that includes ANSI
// escape codes).
func (m *Manager) getReverse(id string) (resolved string) {
	m.idMu.RLock()
	resolved = m.rids[id]
	m.idMu.RUnlock()
	return resolved
}

func (m *Manager) worker() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case xy := <-m.setChan:
			m.mapMu.Lock()
			m.mapping[m.getReverse(xy.id)] = xy
			m.mapMu.Unlock()
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
	vLen := len(v)
	start := -1
	var end, i, w, newlines, newlinesAtStart, lastNewline, width int
	var id string
	var r rune

	for {
		if i+1 >= vLen {
			return v
		}

		r, w = utf8.DecodeRuneInString(v[i:])

		switch r {
		case utf8.RuneError:
			i += w // Skip invalid rune.
			continue
		case '\n':
			if start == -1 {
				lastNewline = i
			}
			i += w
			newlines++
		case identStart:
			start = i
			newlinesAtStart = newlines
			i += w
		case identEnd:
			i += w
			if start == -1 {
				continue
			}
			end = i

			id = v[start+identStartLen : end-identEndLen]

			// calculate the offset here.
			// newlines = countNewlines(v[:start]) // TODO: replace.

			// lastNewline = strings.LastIndex(v[:start], "\n")
			// if lastNewline == -1 {
			// 	lastNewline = 0
			// }

			width = ansi.PrintableRuneWidth(v[lastNewline:start])

			m.setChan <- &Position{id: id, X: width, Y: newlinesAtStart}
			v = v[:start] + v[end:]
			i = start
			vLen = len(v)

			start = -1
		default:
			i += w
		}
	}
}
