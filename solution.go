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
	hs  []int
	res []image.Point
}

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points. The returned point slice is
// only valid until the next call to Solve.
func (sol *Solver) Solve(data []internal.Building) ([]image.Point, error) {
	if len(data) == 0 {
		return nil, nil
	}

	maxx := 0
	for _, b := range data {
		if x := b.Sides[1]; x > maxx {
			maxx = x
		}
	}

	sol.alloc(len(data), maxx+1)

	for _, b := range data {
		for x1, x2 := b.Sides[0], b.Sides[1]; x1 <= x2; x1++ {
			if h := b.Height; sol.hs[x1] < h {
				sol.hs[x1] = h
			}
		}
	}

	return traceHeights(sol.res, sol.hs), nil
}

func traceHeights(res []image.Point, hs []int) []image.Point {
	ch := 0
	x := 0
	for ; x < len(hs); x++ {
		if h := hs[x]; h < ch {
			res = append(res, image.Pt(x-1, ch), image.Pt(x-1, h))
			ch = h
		} else if h > ch {
			res = append(res, image.Pt(x, ch), image.Pt(x, h))
			ch = h
		}
	}
	if ch != 0 {
		res = append(res, image.Pt(x-1, ch), image.Pt(x-1, 0))
	}
	return res
}

func (sol *Solver) alloc(n, hn int) {
	if m := 4 * n; m <= cap(sol.res) {
		sol.res = sol.res[:0]
	} else {
		sol.res = make([]image.Point, 0, m)
	}

	if hn <= cap(sol.hs) {
		sol.hs = sol.hs[:hn]
		for i := range sol.hs {
			sol.hs[i] = 0
		}
	} else {
		sol.hs = make([]int, hn)
	}
}
