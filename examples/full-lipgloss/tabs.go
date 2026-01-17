// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

var (
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		Padding(0, 1)

	activeTab = tab.Border(activeTabBorder, true)

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
)

type tabs struct {
	id     string
	width  int
	active string
	items  []string
}

func (m *tabs) Init() tea.Cmd {
	return nil
}

func (m *tabs) GetHeight() int {
	return lipgloss.Height(m.View())
}

func (m *tabs) Update(msg tea.Msg) tea.Cmd { //nolint:unparam
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.MouseReleaseMsg:
		if msg.Button != tea.MouseLeft {
			return nil
		}

		for _, item := range m.items {
			// Check each item to see if it's in bounds.
			if zone.Get(m.id + item).InBounds(msg) {
				m.active = item
				break
			}
		}

		return nil
	}
	return nil
}

func (m *tabs) View() string {
	out := []string{}

	for _, item := range m.items {
		// Make sure to mark each tab when rendering.
		if item == m.active {
			out = append(out, zone.Mark(m.id+item, activeTab.BorderForeground(highlight).Render(item)))
		} else {
			out = append(out, zone.Mark(
				m.id+item,
				tab.BorderForeground(highlight).Render(item)),
			)
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, out...)
	gap := tabGap.BorderForeground(highlight).Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	return row
}
