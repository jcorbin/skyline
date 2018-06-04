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

	o1 := sol.o1
	x1 := sol.x1
	o2 := sol.o2
	x2 := sol.x2
	h := sol.h
	res := sol.res

	if cap(o1) < len(data) {
		o1 = make([]int, 0, len(data))
		x1 = make([]int, 0, len(data))
		o2 = make([]int, 0, len(data))
		x2 = make([]int, 0, len(data))
		h = make([]int, 0, len(data))
		res = make([]image.Point, 0, 4*len(data))
	} else {
		o1 = o1[:0]
		x1 = x1[:0]
		o2 = o2[:0]
		x2 = x2[:0]
		h = h[:0]
		res = res[:0]
	}

	for i := range data {
		o1, x1 = addOrderedXPoint(o1, x1, i, data[i].Sides[0])
		o2, x2 = addOrderedXPoint(o2, x2, i, data[i].Sides[1])
		h = append(h, data[i].Height)
	}

	sol.o1 = o1
	sol.x1 = x1
	sol.o2 = o2
	sol.x2 = x2
	sol.h = h
	sol.res = res

	cur := image.ZP
	o2i := 0
	for o1i := 0; o1i < len(o1); o1i++ {
		// NOTE test probably doesn't catch edge case where several co-incident
		// opens cause redundant co-linear points
		i := o1[o1i]
		bx := x1[i]

		// close opened past bx
		for o2i < len(o2) && x2[o2[o2i]] < bx {
			j := o2[o2i]
			o2i++
			if rh := remHeightUnder(bx, o2i, o2, x1, x2, h); rh < cur.Y {
				cur, res = tox(cur, res, x2[j])
				cur, res = goy(cur, res, rh)
			}
		}

		// open data[i]
		if bh := h[i]; bh > cur.Y {
			cur, res = tox(cur, res, bx)
			cur, res = goy(cur, res, bh)
		}
	}

	// flush opened
	for o2i < len(o2) {
		j := o2[o2i]
		o2i++
		if rh := remHeight(o2i, o2, x1, x2, h); rh < cur.Y {
			cur, res = tox(cur, res, x2[j])
			cur, res = goy(cur, res, rh)
		}
	}

	return res, nil
}

// TODO pre-compute this rather than this expensive loop every time
func remHeightUnder(bx, o2i int, o2, x1, x2, h []int) int {
	rh := 0
	for ; o2i < len(o2); o2i++ {
		// TODO better to keep an open []int index?
		if k := o2[o2i]; x1[k] < bx {
			if kh := h[k]; kh > rh {
				rh = kh
			}
		}
	}
	return rh
}

// TODO pre-compute this rather than this expensive loop every time
func remHeight(o2i int, o2, x1, x2, h []int) int {
	rh := 0
	for ; o2i < len(o2); o2i++ {
		// TODO better to keep an open []int index?
		if kh := h[o2[o2i]]; kh > rh {
			rh = kh
		}
	}
	return rh
}

func addOrderedXPoint(os, xs []int, i, x int) (_, _ []int) {
	oi := sort.Search(len(os), func(oi int) bool { return xs[os[oi]] > x })
	xs = append(xs, x)
	if oi == len(os) {
		os = append(os, i)
	} else {
		os = append(os, i)
		copy(os[oi+1:], os[oi:])
		os[oi] = i
	}
	return os, xs
}

func goy(cur image.Point, res []image.Point, y int) (image.Point, []image.Point) {
	cur.Y = y
	res = append(res, cur)
	return cur, res
}

func tox(cur image.Point, res []image.Point, x int) (image.Point, []image.Point) {
	if x != cur.X {
		cur.X = x
		res = append(res, cur)
	}
	return cur, res
}
