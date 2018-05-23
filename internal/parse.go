package internal

import (
	"bufio"
	"fmt"
	"image"
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

		n, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return res, fmt.Errorf("invalid x1=%q in %q", parts[0], sc.Text())
		}
		b.Sides[0] = int(n)

		n, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return res, fmt.Errorf("invalid x2=%q in %q", parts[1], sc.Text())
		}
		b.Sides[1] = int(n)

		n, err = strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return res, fmt.Errorf("invalid h=%q in %q", parts[2], sc.Text())
		}
		b.Height = int(n)

		res = append(res, b)
	}

	return res, sc.Err()
}

// ScanPoints reads csv point output; it expects a "x,y" header line, and
// then parses integer pairs from the remaining lines. Any (maybe partial)
// results and error encountered are returned .
func ScanPoints(in io.Reader) ([]image.Point, error) {
	var res []image.Point
	sc := bufio.NewScanner(in)
	if sc.Scan() {
		if line := sc.Text(); line != "x,y" {
			return res, fmt.Errorf("expected point header line, got %q", line)
		}
	}

	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), ",", 2)
		if len(parts) < 2 {
			return res, fmt.Errorf("short line %q", sc.Text())
		}

		var p image.Point

		n, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return res, fmt.Errorf("invalid x=%q in %q", parts[0], sc.Text())
		}
		p.X = int(n)

		n, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return res, fmt.Errorf("invalid y=%q in %q", parts[1], sc.Text())
		}
		p.Y = int(n)

		res = append(res, p)
	}

	return res, sc.Err()
}
