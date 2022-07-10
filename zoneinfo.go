package zone

import tea "github.com/charmbracelet/bubbletea"

// Position is a struct that holds the X and Y coordinates of an zone (start and
// end).
type Position struct {
	id string // raw id of the offset.
	X  int    // X coordinate, starting from 0 (width).
	Y  int    // Y coordinate, starting from 0 (height).
}

// IsZero returns true if the position doesn't reference an offset. Useful when
// calling Get() using an ID that hasn't been registered yet.
func (p *Position) IsZero() bool {
	if p == nil {
		return true
	}
	return p.id == ""
}

type ZoneInfo struct {
	Start *Position
	End   *Position
}

func (z *ZoneInfo) IsZero() bool {
	if z == nil {
		return true
	}
	return z.Start.IsZero() || z.End.IsZero()
}

func (z *ZoneInfo) InBounds(e tea.MouseMsg) bool {
	if z.IsZero() {
		return false
	}

	if z.Start.X > z.End.X || z.Start.Y > z.End.Y {
		return false
	}

	if e.X < z.Start.X || e.Y < z.Start.Y {
		return false
	}

	if e.X >= z.End.X || e.Y > z.End.Y {
		return false
	}

	return true
}

// Offset returns the X and Y offset from the top left of the window, starting at
// 0, 0. Returns -1, -1 if the zone isn't known yet.
func (z *ZoneInfo) Offset(msg tea.MouseMsg) (x, y int) {
	if z.IsZero() {
		return -1, -1
	}

	return msg.X - z.Start.X, msg.Y - z.Start.Y
}
