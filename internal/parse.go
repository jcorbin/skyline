package internal

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Building represents a skyline building in terms of its x-coordinates of its
// sides and its height.
type Building struct {
	Sides  [2]int
	Height int
}

// ScanBuildings reads gen-created csv output; it expects a "x1,x2,h" header
// line, and then parses integer triples from the remaining lines. Any (maybe
// partial) results and error encountered are returned .
func ScanBuildings(in io.Reader) ([]Building, error) {
	var res []Building
	sc := bufio.NewScanner(in)
	if sc.Scan() {
		if line := sc.Text(); line != "x1,x2,h" {
			return res, fmt.Errorf("expected gen header line, got %q", line)
		}
	}

	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), ",", 3)
		if len(parts) < 3 {
			return res, fmt.Errorf("short line %q", sc.Text())
		}

		var b Building

		x1, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return res, fmt.Errorf("invalid x1=%q in %q", parts[0], sc.Text())
		}
		b.Sides[0] = int(x1)

		x2, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return res, fmt.Errorf("invalid x2=%q in %q", parts[1], sc.Text())
		}
		b.Sides[1] = int(x2)

		h, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return res, fmt.Errorf("invalid h=%q in %q", parts[2], sc.Text())
		}
		b.Height = int(h)

		res = append(res, b)
	}

	return res, sc.Err()
}
