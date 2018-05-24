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

	var cur image.Point
	res = append(res, cur)

	for _, b := range data {
		// advance cur.X to left edge
		for len(open) > 0 && b.Sides[0] > open[0].Sides[1] {
			// TODO heap pop
			c := open[0]
			copy(open, open[1:])
			open = open[:len(open)-1]

			cur.X = c.Sides[1]
			res = append(res, cur)
			if h := maxHeight(open); h != cur.Y {
				cur.Y = h
				res = append(res, cur)
			}
		}
		cur.X = b.Sides[0]

		// expand cur.Y to encompass height
		if len(open) == 0 || b.Height > cur.Y {
			res = append(res, cur)
			cur.Y = b.Height
			res = append(res, cur)
		}

		// prune obsoleted buildings
		open = prunePast(open, cur.X)

		// TODO heap insert
		open = append(open, b)
		sort.Slice(open, func(i, j int) bool {
			return open[i].Sides[1] < open[j].Sides[1]
		})
	}

	if len(open) > 0 {
		// TODO if we do end up heap-ing, this would either need to be a
		// heap-pop, or we need to fully sort the remainder before loop
		i := len(open) - 1
		if c := open[i]; c.Sides[1] >= cur.X {
			cur.X = c.Sides[1]
			res = append(res, cur)
			if h := 0; h != cur.Y {
				cur.Y = h
				res = append(res, cur)
			}
			open = prunePast(open, cur.X)
		}
	}

	return res, nil
}

func prunePast(open []internal.Building, x int) []internal.Building {
	for i, j := 0, len(open)-1; j >= 0 && i <= j; {
		if x >= open[i].Sides[1] {
			open[i], open = open[j], open[:j]
			j--
		} else {
			i++
		}
	}
	return open
}

func maxHeight(bs []internal.Building) (h int) {
	for _, b := range bs {
		if b.Height > h {
			h = b.Height
		}
	}
	return h
}
