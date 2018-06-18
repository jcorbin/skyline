package main

import (
	"image"

	"github.com/jcorbin/skyline/internal"
)

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points.
func Solve(data []internal.Building) ([]image.Point, error) {
	var sol Solver
	return sol.Solve(data)
}

// Solver holds any state for solving the skyline problem, potentially re-using
// previously allocated state memory.
type Solver struct {
	res []image.Point
}

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points. The returned point slice is
// only valid until the next call to Solve.
func (sol *Solver) Solve(data []internal.Building) ([]image.Point, error) {
	if len(data) == 0 {
		return nil, nil
	}

	minx := data[0].Sides[0]
	maxx := data[0].Sides[1]
	maxh := data[0].Height
	for i := 1; i < len(data); i++ {
		if x := data[i].Sides[0]; minx > x {
			minx = x
		}
		if x := data[i].Sides[1]; maxx < x {
			maxx = x
		}
		if h := data[i].Height; maxh < h {
			maxh = h
		}
	}

	hs := make([]int, (maxx-minx)+1)

	for h := 1; h <= maxh; h++ {
		for i := 0; i < len(data); i++ {
			if data[i].Height == h {
				for x1, x2 := data[i].Sides[0]-minx, data[i].Sides[1]-minx; x1 <= x2; x1++ {
					hs[x1] = h
				}
			}
		}
	}

	if m := 4 * len(data); m > cap(sol.res) {
		sol.res = make([]image.Point, 0, m)
	}

	res := sol.res
	ch := 0
	x := minx
	for i := 0; i < len(hs); i++ {
		if h := hs[i]; h < ch {
			res = append(res, image.Pt(x-1, ch), image.Pt(x-1, h))
			ch = h
		} else if h > ch {
			res = append(res, image.Pt(x, ch), image.Pt(x, h))
			ch = h
		}
		x++
	}
	if ch != 0 {
		res = append(res, image.Pt(x-1, ch), image.Pt(x-1, 0))
	}

	return res, nil
}
