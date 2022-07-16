// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1).
			MarginRight(2)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)
)

type dialog struct {
	id     string
	height int
	width  int

	active   string
	question string
}

func (m dialog) Init() tea.Cmd {
	return nil
}

func (m dialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			return m, nil
		}

		if zone.Get(m.id + "confirm").InBounds(msg) {
			m.active = "confirm"
		} else if zone.Get(m.id + "cancel").InBounds(msg) {
			m.active = "cancel"
		}

		return m, nil
	}
	return m, nil
}

func (m dialog) View() string {
	var okButton, cancelButton string

	if m.active == "confirm" {
		okButton = activeButtonStyle.Render("Yes")
		cancelButton = buttonStyle.Render("Maybe")
	} else {
		okButton = buttonStyle.Render("Yes")
		cancelButton = activeButtonStyle.Render("Maybe")
	}

	question := lipgloss.NewStyle().Width(27).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		zone.Mark(m.id+"confirm", okButton),
		zone.Mark(m.id+"cancel", cancelButton),
	)
	return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, buttons))
}
