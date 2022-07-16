// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type history struct {
	id     string
	height int
	width  int

	active string
	items  []string
}

func (m history) Init() tea.Cmd {
	return nil
}

func (m history) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			return m, nil
		}

		for _, item := range m.items {
			// Check each item to see if it's in bounds.
			if zone.Get(m.id + item).InBounds(msg) {
				m.active = item
				break
			}
		}
	}
	return m, nil
}

func (m history) View() string {
	historyStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(subtle).
		Margin(1).
		Padding(1, 2).
		Width((m.width / len(m.items)) - 2).
		Height(m.height - 2).
		MaxHeight(m.height)

	out := []string{}

	for _, item := range m.items {
		if item == m.active {
			// Customize the active item.
			out = append(out, zone.Mark(m.id+item, historyStyle.Copy().Background(highlight).Render(item)))
		} else {
			// Make sure to mark all zones.
			out = append(out, zone.Mark(m.id+item, historyStyle.Render(item)))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, out...)
}
