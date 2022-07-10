// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package zone

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func TestMain(m *testing.M) {
	NewGlobal(500)

	testsScan = []scanTestCase{
		{"empty", "", "", nil},
		{"single", "a", "a", nil},
		{"double", "aa", "aa", nil},
		{"triple", "aaa", "aaa", nil},
		{"quad", "aaaa", "aaaa", nil},
		{"lipgloss-empty", testStyle.Render(""), testStyle.Render(""), nil},
		{"lipgloss-basic", testStyle.Render("testing"), testStyle.Render("testing"), nil},
		{"lipgloss-basic-start", "a" + testStyle.Render("testing"), "a" + testStyle.Render("testing"), nil},
		{"lipgloss-basic-end", testStyle.Render("testing") + "a", testStyle.Render("testing") + "a", nil},
		{"lipgloss-basic-start-end", "a" + testStyle.Render("testing") + "a", "a" + testStyle.Render("testing") + "a", nil},
		{"lipgloss-basic-between", testStyle.Render("testing") + "a" + testStyle.Render("testing"), testStyle.Render("testing") + "a" + testStyle.Render("testing"), nil},
		{"id-empty", Mark("testing", ""), "", []string{"testing"}},
		{"id-single-start", "a" + Mark("testing", "a"), "aa", []string{"testing"}},
		{"id-single-end", Mark("testing", "a") + "a", "aa", []string{"testing"}},
		{"id-single-start-end", "a" + Mark("testing", "b") + "a", "aba", []string{"testing"}},
		{"id-single-between", Mark("testing", "b") + "a" + Mark("testing", "b"), "bab", []string{"testing"}},
		{"id-with-lipgloss-start", testStyle.Render(Mark("testing", "testing") + "testing"), testStyle.Render("testingtesting"), []string{"testing"}},
		{"id-with-lipgloss-end", testStyle.Render("testing" + Mark("testing", "testing")), testStyle.Render("testingtesting"), []string{"testing"}},
		{"id-multi-empty", Mark("foo", "") + Mark("bar", ""), "", []string{"foo", "bar"}},
		{"id-multi-start", "a" + Mark("foo", "") + Mark("bar", ""), "a", []string{"foo", "bar"}},
		{"id-multi-end", Mark("foo", "") + Mark("bar", "") + "a", "a", []string{"foo", "bar"}},
		{"id-multi-start-end", "a" + Mark("foo", "") + Mark("bar", "") + "a", "aa", []string{"foo", "bar"}},
		{"long-x1", "a" + Mark("longtest", longStyle.Render("testing")) + "a", "a" + longStyle.Render("testing") + "a", []string{"longtest"}},
		{"long-x2", strings.Repeat("a"+Mark("longtest", longStyle.Render("testing"))+"a", 1), strings.Repeat("a"+longStyle.Render("testing")+"a", 1), []string{"longtest"}},
		{"long-x4", strings.Repeat("a"+Mark("longtest", longStyle.Render("testing"))+"a", 4), strings.Repeat("a"+longStyle.Render("testing")+"a", 4), []string{"longtest"}},
		{"long-x6", strings.Repeat("a"+Mark("longtest", longStyle.Render("testing"))+"a", 6), strings.Repeat("a"+longStyle.Render("testing")+"a", 6), []string{"longtest"}},
		{"long-x8", strings.Repeat("a"+Mark("longtest", longStyle.Render("testing"))+"a", 8), strings.Repeat("a"+longStyle.Render("testing")+"a", 8), []string{"longtest"}},
		{"long-x10", strings.Repeat("a"+Mark("longtest", longStyle.Render("testing"))+"a", 10), strings.Repeat("a"+longStyle.Render("testing")+"a", 10), []string{"longtest"}},
	}

	m.Run()
	Close()
}

type scanTestCase struct {
	name string
	in   string
	want string
	ids  []string
}

var (
	testStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#383838")).
			Bold(true).
			Italic(true).
			Blink(true)

	longStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#383838")).
			Bold(true).
			Italic(true).
			Blink(true).
			Underline(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F12356")).
			BorderBackground(lipgloss.Color("#459082")).
			Padding(5, 4)
	testsScan []scanTestCase
)

func BenchmarkScan(b *testing.B) {
	for _, test := range testsScan {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Scan(test.in)
			}
		})
	}
}

func TestScan(t *testing.T) {
	for _, test := range testsScan {
		t.Run(test.name, func(t *testing.T) {
			got := Scan(test.in)
			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
			if len(test.ids) > 0 {
				time.Sleep(15 * time.Millisecond)
				for _, id := range test.ids {
					if xy := Get(id); xy.IsZero() {
						t.Errorf("id %q not found", id)
					}
				}
			}
		})
	}
}

func FuzzScan(f *testing.F) {
	for _, test := range testsScan {
		f.Add(test.in)
		f.Add(test.want)
	}

	f.Fuzz(func(t *testing.T, a string) {
		_ = Scan(a)
	})
}

func BenchmarkMark(b *testing.B) {
	for _, test := range testsScan {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Mark(test.name, test.in)
			}
		})
	}
}
