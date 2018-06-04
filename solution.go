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
	o1  []int
	x1  []int
	o2  []int
	x2  []int
	h   []int
	res []image.Point
}

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points. The returned point slice is
// only valid until the next call to Solve.
func (sol *Solver) Solve(data []internal.Building) ([]image.Point, error) {
	if len(data) == 0 {
		return nil, nil
	}

	if cap(sol.o1) < len(data) {
		sol.o1 = make([]int, 0, len(data))
		sol.x1 = make([]int, 0, len(data))
		sol.o2 = make([]int, 0, len(data))
		sol.x2 = make([]int, 0, len(data))
		sol.h = make([]int, 0, len(data))
		sol.res = make([]image.Point, 0, 4*len(data))
	} else {
		sol.o1 = sol.o1[:0]
		sol.x1 = sol.x1[:0]
		sol.o2 = sol.o2[:0]
		sol.x2 = sol.x2[:0]
		sol.h = sol.h[:0]
		sol.res = sol.res[:0]
	}

	for i := range data {
		sol.o1, sol.x1 = addOrderedXPoint(sol.o1, sol.x1, i, data[i].Sides[0])
		sol.o2, sol.x2 = addOrderedXPoint(sol.o2, sol.x2, i, data[i].Sides[1])
		sol.h = append(sol.h, data[i].Height)
	}

	cur := image.ZP
	o2i := 0
	for o1i := 0; o1i < len(sol.o1); o1i++ {
		// NOTE test probably doesn't catch edge case where several co-incident
		// opens cause redundant co-linear points
		i := sol.o1[o1i]
		bx := sol.x1[i]

		// close opened past bx
		for o2i < len(sol.o2) && sol.x2[sol.o2[o2i]] < bx {
			j := sol.o2[o2i]
			cur, sol.res = tox(cur, sol.res, sol.x2[j])

			// TODO try to pre-compute
			o2i++
			rh := 0
			for k := o2i; k < len(sol.o2) && sol.x2[sol.o2[k]] < bx; k++ {
				if sol.h[k] > rh {
					rh = sol.h[k]
				}
			}

			cur, sol.res = downy(cur, sol.res, rh)
		}

		// open data[i]
		cur, sol.res = tox(cur, sol.res, bx)
		cur, sol.res = yup(cur, sol.res, sol.h[i])
	}

	// flush opened
	for o2i < len(sol.o2) {
		j := sol.o2[o2i]
		cur, sol.res = tox(cur, sol.res, sol.x2[j])

		// TODO try to pre-compute
		o2i++
		rh := 0
		for k := o2i; k < len(sol.o2); k++ {
			if sol.h[k] > rh {
				rh = sol.h[k]
			}
		}

		cur, sol.res = downy(cur, sol.res, rh)
	}

	return sol.res, nil
}

func addOrderedXPoint(os, xs []int, i, x int) (_, _ []int) {
	oi := sort.Search(len(os), func(oi int) bool { return xs[os[oi]] > x })
	xs = append(xs, x)
	if oi == len(os) {
		os = append(os, i)
	} else {
		os = append(os, i)
		copy(os[oi+1:], os[oi:])
	}
	return os, xs
}

func tox(cur image.Point, res []image.Point, x int) (image.Point, []image.Point) {
	if x != cur.X {
		cur.X = x
		res = append(res, cur)
	}
	return cur, res
}

func yup(cur image.Point, res []image.Point, y int) (image.Point, []image.Point) {
	if y > cur.Y {
		cur.Y = y
		res = append(res, cur)
	}
	return cur, res
}

func downy(cur image.Point, res []image.Point, y int) (image.Point, []image.Point) {
	if y < cur.Y {
		cur.Y = y
		res = append(res, cur)
	}
	return cur, res
}
