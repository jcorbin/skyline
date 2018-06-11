package main

import (
	"image"
	"sort"

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
	x1 []int
	x2 []int
	h  []int

	cur image.Point
	res []image.Point
}

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points. The returned point slice is
// only valid until the next call to Solve.
func (sol *Solver) Solve(data []internal.Building) ([]image.Point, error) {
	sol.alloc(data)
	sol.cur = image.ZP

	for _, b := range data {
		x1, x2, h := b.Sides[0], b.Sides[1], b.Height
		n := len(sol.x2)
		i := sort.Search(n, func(i int) bool { return sol.x2[i] > x1 })
		j := sort.Search(n, func(i int) bool { return sol.x2[i] > x2 })
		keep := make([]bool, j-i+1)
		// TODO

	}

	for i := range sol.x1 {
		sol.tox(sol.x1[i])
		sol.toy(sol.h[i])
		sol.gox(sol.x2[i])
		sol.toy(0)
	}

	return sol.res, nil
}

func (sol *Solver) alloc(data []internal.Building) {
	n := len(data)
	if n <= cap(sol.x1) {
		sol.x1 = sol.x1[:0]
		sol.x2 = sol.x2[:0]
		sol.h = sol.h[:0]
		sol.res = sol.res[:0]
	} else {
		sol.x1 = make([]int, 0, n)
		sol.x2 = make([]int, 0, n)
		sol.h = make([]int, 0, n)
		sol.res = make([]image.Point, 0, 4*n)
	}
}

func (sol *Solver) tox(x int) {
	if x != sol.cur.X {
		sol.gox(x)
	}
}
func (sol *Solver) toy(y int) {
	if y != sol.cur.Y {
		sol.goy(y)
	}
}
func (sol *Solver) gox(x int)
func (sol *Solver) goy(y int)
