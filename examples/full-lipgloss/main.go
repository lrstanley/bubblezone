// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
	zone "github.com/lrstanley/bubblezone/v2"
)

// This is a modified version of this example, supporting full screen, dynamic
// resizing, and clickable models (tabs, lists, dialogs, etc).
// 	https://github.com/charmbracelet/lipgloss/blob/master/example

var (
	subtle    = compat.AdaptiveColor{Light: lipgloss.Color("#D9DCCF"), Dark: lipgloss.Color("#383838")}
	highlight = compat.AdaptiveColor{Light: lipgloss.Color("#874BFD"), Dark: lipgloss.Color("#7D56F4")}
	special   = compat.AdaptiveColor{Light: lipgloss.Color("#43BF6D"), Dark: lipgloss.Color("#73F59F")}
)

type model struct {
	height int
	width  int

	tabs    tea.Model
	dialog  tea.Model
	list1   tea.Model
	list2   tea.Model
	history tea.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) isInitialized() bool {
	return m.height != 0 && m.width != 0
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.isInitialized() {
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Example of toggling mouse event tracking on/off.
		if msg.String() == "ctrl+e" {
			zone.SetEnabled(!zone.Enabled())
			return m, nil
		}

		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		msg.Height -= 2
		msg.Width -= 4
		return m.propagate(msg), nil
	}

	return m.propagate(msg), nil
}

func (m *model) propagate(msg tea.Msg) tea.Model {
	// Propagate to all children.
	m.tabs, _ = m.tabs.Update(msg)
	m.dialog, _ = m.dialog.Update(msg)
	m.list1, _ = m.list1.Update(msg)
	m.list2, _ = m.list2.Update(msg)

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		msg.Height -= m.tabs.(tabs).height + m.list1.(list).height
		m.history, _ = m.history.Update(msg)
		return m
	}

	m.history, _ = m.history.Update(msg)
	return m
}

func (m model) View() string {
	if !m.isInitialized() {
		return ""
	}

	s := lipgloss.NewStyle().MaxHeight(m.height).MaxWidth(m.width).Padding(1, 2, 1, 2)

	return zone.Scan(s.Render(lipgloss.JoinVertical(lipgloss.Top,
		m.tabs.(tea.ViewModel).View(), "",
		lipgloss.PlaceHorizontal(
			m.width, lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.list1.(tea.ViewModel).View(), m.list2.(tea.ViewModel).View(), m.dialog.(tea.ViewModel).View(),
			),
			lipgloss.WithWhitespaceChars(" "),
		),
		m.history.(tea.ViewModel).View(),
	)))
}

func main() {
	// Initialize a global zone manager, so we don't have to pass around the manager
	// throughout components.
	zone.NewGlobal()

	m := &model{
		tabs: &tabs{
			id:     zone.NewPrefix(), // Give each type an ID, so no zones will conflict.
			height: 3,
			active: "Lip Gloss",
			items:  []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"},
		},
		dialog: &dialog{
			id:       zone.NewPrefix(),
			height:   8,
			active:   "confirm",
			question: "Are you sure you want to eat marmalade?",
		},
		list1: &list{
			id:     zone.NewPrefix(),
			height: 8,
			title:  "Citrus Fruits to Try",
			items: []listItem{
				{name: "Grapefruit", done: true},
				{name: "Yuzu", done: false},
				{name: "Citron", done: false},
				{name: "Kumquat", done: true},
				{name: "Pomelo", done: false},
			},
		},
		list2: &list{
			id:     zone.NewPrefix(),
			height: 8,
			title:  "Actual Lip Gloss Vendors",
			items: []listItem{
				{name: "Glossier", done: true},
				{name: "Claire's Boutique", done: true},
				{name: "Nyx", done: false},
				{name: "Mac", done: false},
				{name: "Milk", done: false},
			},
		},
		history: &history{
			id: zone.NewPrefix(),
			items: []string{
				"The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum: Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos.",
				"Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac.",
				"In 1524, Henry VIII, King of England, received a “box of marmalade” from Mr. Hull of Exeter. This was probably marmelada, a solid quince paste from Portugal, still made and sold in southern Europe today. It became a favourite treat of Anne Boleyn and her ladies in waiting.",
			},
		},
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}
