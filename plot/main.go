package main

import (
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jcorbin/skyline/internal"
)

func main() {
	if err := run(os.Stdin); err != nil {
		log.Fatalln(err)
	}
}

func run(in io.Reader) error {
	points, err := internal.ScanPoints(in)
	if err != nil {
		return err
	}

	// measure world bounding box
	var size image.Point
	for _, p := range points {
		if p.X > size.X {
			size.X = p.X
		}
		if p.Y > size.Y {
			size.Y = p.Y
		}
	}

	// inflate bounding box for reasons
	if dw := size.X / 4; dw > 2 {
		size.X += dw
	} else {
		size.X += 2
	}
	if dh := size.Y / 4; dh > 2 {
		size.Y += dh
	} else {
		size.Y += 2
	}

	// draw points into an output array; display screen cells with double
	// width to improve terminal aspect ratio.
	size.X *= 2
	out := make([]byte, size.X*size.Y)
	var cur image.Point
	for _, p := range points {
		if p.Eq(cur) {
			continue
		}

		if p.X < cur.X {
			return fmt.Errorf("unsupported backwards X scan")
		}

		for p.X > cur.X {
			out[cur.Y*size.X+2*cur.X] = '-'
			out[cur.Y*size.X+2*cur.X+1] = '-'
			cur.X++
		}

		out[cur.Y*size.X+2*cur.X] = '['
		out[cur.Y*size.X+2*cur.X+1] = ']'

		for p.Y > cur.Y {
			cur.Y++
			out[cur.Y*size.X+2*cur.X] = '|'
			// out[cur.Y*size.X+2*cur.X+1] = '|'
		}

		for p.Y < cur.Y {
			cur.Y--
			// out[cur.Y*size.X+2*cur.X] = '|'
			out[cur.Y*size.X+2*cur.X+1] = '|'
		}
	}

	// render the output array, each cell to a 2x1 pair of screen cells
	fmt.Printf("/%s\\\n", strings.Repeat("=", size.X))
	for y := size.Y - 1; y >= 0; y-- {
		i := y * size.X
		buf := make([]byte, size.X)
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
	fmt.Printf("\\%s/\n", strings.Repeat("=", size.X))

	return nil
}
