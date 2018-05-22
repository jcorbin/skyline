package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := run(os.Stdin); err != nil {
		log.Fatalln(err)
	}
}

func run(in io.Reader) error {
	data, err := ScanBuildings(in)
	if err != nil {
		return err
	}

	// measure world bounding box
	w, h := 0, 0
	for _, b := range data {
		if b.Sides[1] > w {
			w = b.Sides[1]
		}
		if b.Height > h {
			h = b.Height
		}
	}
	if dw := w / 4; dw > 2 {
		w += dw
	} else {
		w += 2
	}
	if dh := h / 4; dh > 2 {
		h += dh
	} else {
		h += 2
	}

	// output array to lay building out in; double width to improve terminal
	// aspect ratio.
	out := make([]byte, 2*w*h)
	for _, b := range data {
		x := b.Sides[0]
		y := 0
		for ; y < b.Height; y++ {
			out[y*2*w+2*x] = '|'
		}

		out[y*2*w+2*x] = '/'
		out[y*2*w+2*x+1] = '-'
		x++

		for ; x < b.Sides[1]; x++ {
			out[y*2*w+2*x] = '-'
			out[y*2*w+2*x+1] = '-'
		}

		out[y*2*w+2*x] = '-'
		out[y*2*w+2*x+1] = '\\'
		y--

		for ; y >= 0; y-- {
			out[y*2*w+2*x+1] = '|'
		}
	}

	// render the output array, each cell to a 2x1 pair of screen cells
	fmt.Printf("/%s\\\n", strings.Repeat("=", 2*w))
	for y := h - 1; y >= 0; y-- {
		i := y * 2 * w
		buf := make([]byte, 2*w)
		for j := 0; j < len(buf); {
			b := out[i]
			if b == 0 {
				b = ' '
			}
			buf[j] = b
			j++
			i++
		}
		fmt.Printf("|%s|\n", buf)
	}
	fmt.Printf("\\%s/\n", strings.Repeat("=", 2*w))

	return nil
}

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
