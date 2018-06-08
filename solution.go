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
	cur image.Point
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

	sol.alloc(len(data), maxx)

	for _, b := range data {
		x1, x2, h := b.Sides[0], b.Sides[1], b.Height
		for x := x1; x <= x2; x++ {
			if sol.hs[x] < h {
				sol.hs[x] = h
			}
		}
	}

	x := 0
	for ; x <= maxx; x++ {
		if h := sol.hs[x]; h < sol.cur.Y {
			sol.gox(x - 1)
			sol.goy(h)
		} else if h > sol.cur.Y {
			sol.gox(x)
			sol.goy(h)
		}
	}
	if sol.cur.Y != 0 {
		sol.gox(maxx)
		sol.goy(0)
	}

	return sol.res, nil
}

func (sol *Solver) gox(x int) {
	sol.cur.X = x
	sol.res = append(sol.res, sol.cur)
}

func (sol *Solver) goy(y int) {
	sol.cur.Y = y
	sol.res = append(sol.res, sol.cur)
}

func (sol *Solver) alloc(n, maxx int) {
	if m := 4 * n; m <= cap(sol.res) {
		sol.res = sol.res[:0]
	} else {
		sol.res = make([]image.Point, 0, m)
	}

	if hn := maxx + 1; hn <= cap(sol.hs) {
		sol.hs = sol.hs[:hn]
		for i := range sol.hs {
			sol.hs[i] = 0
		}
	} else {
		sol.hs = make([]int, hn)
	}
}
