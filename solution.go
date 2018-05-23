package main

import (
	"image"
	"sort"

	"github.com/jcorbin/skyline/internal"
)

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points.
func Solve(data []internal.Building) ([]image.Point, error) {
	if len(data) == 0 {
		return nil, nil
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Sides[0] < data[j].Sides[0]
	})

	res := make([]image.Point, 0, 4*len(data))
	open := make([]internal.Building, 0, len(data))
	openHeight := 0

	var cur image.Point
	res = append(res, cur)

	for _, b := range data {
		for len(open) > 0 && b.Sides[0] > open[0].Sides[1] {
			// TODO heap pop
			c := open[0]
			copy(open, open[1:])
			open = open[:len(open)-1]

			res = append(res, image.Pt(c.Sides[1], openHeight))
			h := maxHeight(open)
			if h != openHeight {
				openHeight = h
				res = append(res, image.Pt(c.Sides[1], openHeight))
			}
		}

		if len(open) == 0 || b.Height > openHeight {
			res = append(res, image.Pt(b.Sides[0], openHeight))
			openHeight = b.Height
			res = append(res, image.Pt(b.Sides[0], openHeight))
		}

		// TODO heap insert
		open = append(open, b)
		sort.Slice(open, func(i, j int) bool {
			return open[i].Sides[1] < open[j].Sides[1]
		})

	}

	for i := len(open); i > 0; {
		i--

		// TODO if we do end up heap-ing, this would either need to be a
		// heap-pop, or we need to fully sort the remainder before loop
		c := open[i]

		openHeight = maxHeight(open[:i])
		res = append(res, image.Pt(c.Sides[1], openHeight))
	}

	return res, nil
}

func maxHeight(bs []internal.Building) (h int) {
	for _, b := range bs {
		if b.Height > h {
			h = b.Height
		}
	}
	return h
}
