// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

var (
	listStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			MarginRight(2)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			MarginRight(2)

	listItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	checkMark = lipgloss.NewStyle().
			SetString("✓").
			PaddingRight(1)

	listDoneStyle = lipgloss.NewStyle().
			Strikethrough(true)
)

type listItem struct {
	name string
	done bool
}

type list struct {
	id    string
	title string
	items []listItem
}

func (m *list) Init() tea.Cmd {
	return nil
}

func (m *list) GetHeight() int {
	return lipgloss.Height(m.View())
}

func (m *list) Update(msg tea.Msg) tea.Cmd { //nolint:unparam
	switch msg := msg.(type) {
	case tea.MouseReleaseMsg:
		if msg.Button != tea.MouseLeft {
			return nil
		}

		for i, item := range m.items {
			// Check each item to see if it's in bounds.
			if zone.Get(m.id + item.name).InBounds(msg) {
				m.items[i].done = !m.items[i].done
				break
			}
		}

		return nil
	}
	return nil
}

func (m *list) View() string {
	out := []string{listHeader.BorderForeground(subtle).Render(m.title)}

	for _, item := range m.items {
		if item.done {
			out = append(out, zone.Mark(
				m.id+item.name,
				checkMark.Foreground(special).String()+
					listDoneStyle.Foreground(completed).Render(item.name),
			))
			continue
		}

		out = append(out, zone.Mark(m.id+item.name, listItemStyle.Render(item.name)))
	}

	return listStyle.BorderForeground(subtle).Render(
		lipgloss.JoinVertical(lipgloss.Left, out...),
	)
}
