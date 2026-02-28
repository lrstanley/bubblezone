// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// This is a modified version of this example, supporting full screen, dynamic
// resizing, and clickable models (tabs, lists, dialogs, etc).
// 	https://github.com/charmbracelet/lipgloss/blob/master/example

var (
	subtle    = lipgloss.Color("#383838")
	highlight = lipgloss.Color("#7D56F4")
	special   = lipgloss.Color("#73F59F")
	completed = lipgloss.Color("#696969")
)

type model struct {
	height int
	width  int

	tabs    *tabs
	dialog  *dialog
	list1   *list
	list2   *list
	history *history
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) isInitialized() bool {
	return m.height != 0 && m.width != 0
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	}

	return m, m.propagate(msg) //nolint:gocritic
}

func (m model) propagate(msg tea.Msg) tea.Cmd {
	// Propagate to all children.
	cmds := []tea.Cmd{
		m.tabs.Update(msg),
		m.dialog.Update(msg),
		m.list1.Update(msg),
		m.list2.Update(msg),
	}

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		msg.Height -= m.tabs.GetHeight() +
			max(m.list1.GetHeight(), m.list2.GetHeight(), m.dialog.GetHeight()) +
			2 // +1 for bottom margin on tabs, +1 for top margin on history.

		cmds = append(cmds, m.history.Update(msg))
		return tea.Batch(cmds...)
	}
	return tea.Batch(append(cmds, m.history.Update(msg))...)
}

func (m model) View() tea.View {
	var view tea.View
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	if !m.isInitialized() {
		return view
	}

	s := lipgloss.NewStyle().MaxHeight(m.height).MaxWidth(m.width)

	// Wrap the main models view in [zone.Scan].
	view.SetContent(zone.Scan(s.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.NewStyle().MarginBottom(1).Render(m.tabs.View()),
			lipgloss.PlaceHorizontal(
				m.width, lipgloss.Center,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					m.list1.View(), m.list2.View(), m.dialog.View(),
				),
				lipgloss.WithWhitespaceChars(" "),
			),
			lipgloss.NewStyle().MarginTop(1).Render(m.history.View()),
		),
	)))
	return view
}

func main() {
	// Initialize a global zone manager, so we don't have to pass around the manager
	// throughout components.
	zone.NewGlobal()

	m := &model{
		tabs: &tabs{
			id:     zone.NewPrefix(), // Give each type an ID, so no zones will conflict.
			active: "Lip Gloss",
			items:  []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"},
		},
		dialog: &dialog{
			id:       zone.NewPrefix(),
			active:   "confirm",
			question: "Are you sure you want to eat marmalade?",
		},
		list1: &list{
			id:    zone.NewPrefix(),
			title: "Citrus Fruits to Try",
			items: []listItem{
				{name: "Grapefruit", done: true},
				{name: "Yuzu", done: false},
				{name: "Citron", done: false},
				{name: "Kumquat", done: true},
				{name: "Pomelo", done: false},
			},
		},
		list2: &list{
			id:    zone.NewPrefix(),
			title: "Actual Lip Gloss Vendors",
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

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("error running program:", err) //nolint:forbidigo
		os.Exit(1)
	}
}
