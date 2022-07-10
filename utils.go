// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package zone

func countNewlines(v string) int {
	var n int
	for _, r := range v {
		if r == '\n' {
			n++
		}
	}

	return n
}

func isNumber(rid string) bool {
	for _, r := range rid {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
