package main

import (
	"fmt"
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
	data, err := internal.ScanBuildings(in)
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
