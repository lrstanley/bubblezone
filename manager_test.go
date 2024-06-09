// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package zone

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func TestMain(m *testing.M) {
	NewGlobal()

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
		{"id-empty", Mark("testing1", ""), "", nil},
		{"id-single-start", "a" + Mark("testing2", "a"), "aa", []string{"testing2"}},
		{"id-single-end", Mark("testing3", "a") + "a", "aa", []string{"testing3"}},
		{"id-single-start-end", "a" + Mark("testing4", "b") + "a", "aba", []string{"testing4"}},
		{"id-single-between", Mark("testing5", "b") + "a" + Mark("testing6", "b"), "bab", []string{"testing5", "testing6"}},
		{"id-with-lipgloss-start", testStyle.Render(Mark("testing7", "testing") + "testing"), testStyle.Render("testingtesting"), []string{"testing7"}},
		{"id-with-lipgloss-end", testStyle.Render("testing" + Mark("testing8", "testing")), testStyle.Render("testingtesting"), []string{"testing8"}},
		{"id-multi-empty", Mark("foo1", "") + Mark("bar1", ""), "", nil},
		{"id-multi-start", "a" + Mark("foo2", "b") + Mark("bar2", "c"), "abc", []string{"foo2", "bar2"}},
		{"id-multi-end", Mark("foo3", "a") + Mark("bar3", "b") + "c", "abc", []string{"foo3", "bar3"}},
		{"id-multi-start-end", "a" + Mark("foo4", "b") + Mark("bar4", "c") + "d", "abcd", []string{"foo4", "bar4"}},
		{"inception", Mark("foo", Mark("bar", "b")), "b", []string{"foo", "bar"}},
		{"long-x1", "a" + Mark("longtest5", longStyle.Render("testing")) + "a", "a" + longStyle.Render("testing") + "a", []string{"longtest5"}},
		{"long-x2", strings.Repeat("a"+Mark("longtest", longStyle.Render("testing"))+"a", 1), strings.Repeat("a"+longStyle.Render("testing")+"a", 1), []string{"longtest"}},
		{"long-x4", strings.Repeat("a"+Mark("longtest", longStyle.Render("testing"))+"a", 4), strings.Repeat("a"+longStyle.Render("testing")+"a", 4), []string{"longtest"}},
		{"long-x6", strings.Repeat("a"+Mark("longtest", longStyle.Render("testing"))+"a", 6), strings.Repeat("a"+longStyle.Render("testing")+"a", 6), []string{"longtest"}},
		{"long-x8", strings.Repeat("a"+Mark("longtest", longStyle.Render("testing"))+"a", 8), strings.Repeat("a"+longStyle.Render("testing")+"a", 8), []string{"longtest"}},
		{"long-x10", strings.Repeat("a"+Mark("longtest", longStyle.Render("testing"))+"a", 10), strings.Repeat("a"+longStyle.Render("testing")+"a", 10), []string{"longtest"}},
		{"invalid-no-bracket", "a\x1B12345Zb", "a\x1B12345Zb", nil},
		{"invalid-no-bracket-end", "a\x1B", "a\x1B", nil},
		{"invalid-no-numbers", "a\x1BZb", "a\x1BZb", nil},
		{"invalid-no-numbers-end", "a\x1BZ", "a\x1BZ", nil},
		{"invalid-marker-end", "a\x1B12345b", "a\x1B12345b", nil},
		{"invalid-marker-end-2", "a\x1B12345", "a\x1B12345", nil},
		{"invalid-run-of-numbers", "a\x1B12345b6Z", "a\x1B12345b6Z", nil},
		{"invalid-misc", "\x1Ba\x1B\x1B\x1B12345b6Z\x1B", "\x1Ba\x1B\x1B\x1B12345b6Z\x1B", nil},
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

	f.Fuzz(func(_ *testing.T, a string) {
		_ = Scan(a)
	})
}

func TestScanDisabled(t *testing.T) {
	zm := New()
	defer zm.Close()

	zm.SetEnabled(false)

	for _, test := range testsScan {
		t.Run(test.name, func(t *testing.T) {
			got := zm.Scan(test.in)
			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestMark(_ *testing.T) {
	var out string
	for _, test := range testsScan[0:10] {
		got := Mark(test.name, test.in)
		out += got
	}

	_ = Scan(out)
	time.Sleep(100 * time.Millisecond)

	for _, test := range testsScan[0:10] {
		_ = Get(test.name)
	}
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

func TestMarkDisabled(t *testing.T) {
	zm := New()
	defer zm.Close()

	zm.SetEnabled(false)

	for _, test := range testsScan[0:10] {
		if got := zm.Mark(test.name, test.in); got != test.in {
			t.Errorf("got %q, want %q", got, test.in)
		}
	}
}

func TestWorkerClear(t *testing.T) {
	_ = Scan("a" + Mark("foo", "b") + "c")
	_ = Scan("a" + Mark("bar", "b") + "c")
	time.Sleep(100 * time.Millisecond)

	if xy := Get("foo"); !xy.IsZero() {
		t.Errorf("%#v not cleared (after %#v)", xy, Get("bar"))
	}
}

func TestClear(t *testing.T) {
	_ = Scan("a" + Mark("foo", "b") + "c")
	Clear("foo")
	if xy := Get("foo"); !xy.IsZero() {
		t.Errorf("%#v not cleared", xy)
	}
}

func TestClose(t *testing.T) {
	mgr := New()
	mgr.Close()
	time.Sleep(100 * time.Millisecond)
	_ = mgr.Scan("a" + Mark("foo", "b") + "c")
	time.Sleep(100 * time.Millisecond)
	if xy := mgr.Get("foo"); !xy.IsZero() {
		t.Errorf("%#v fetched, but closed", xy)
	}
}

func TestGlobalInitialize(_ *testing.T) {
	NewGlobal()
	NewGlobal()
}
