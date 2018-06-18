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
	hs  []int
	dh  []int
	ix  []int
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
	for i := 1; i < len(data); i++ {
		if x := data[i].Sides[0]; minx > x {
			minx = x
		}
		if x := data[i].Sides[1]; maxx < x {
			maxx = x
		}
	}

	hn := (maxx - minx) + 1
	if hn > cap(sol.hs) {
		sol.hs = make([]int, 0, hn)
	}
	hs := sol.hs[:hn]

	if len(data) > cap(sol.ix) {
		sol.ix = make([]int, 0, len(data))
		sol.dh = make([]int, 0, len(data))
	}
	ix := sol.ix[:len(data)]
	dh := sol.dh[:len(data)]

	for bh := heapifyByHeight(data, dh, ix); len(bh.ix) > 0; {
		h, i := bh.pop()
		for x1, x2 := data[i].Sides[0]-minx, data[i].Sides[1]-minx; x1 <= x2; x1++ {
			hs[x1] = h
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

type byH struct {
	dh, ix []int
}

func heapifyByHeight(data []internal.Building, dh, ix []int) byH {
	bh := byH{
		dh: dh[:len(data)],
		ix: ix[:len(data)],
	}
	for i := 0; i < len(data); i++ {
		bh.dh[i] = data[i].Height
		bh.ix[i] = i
	}
	n := len(bh.ix)
	for i := n/2 - 1; i >= 0; i-- {
		bh.down(i, n)
	}
	return bh
}

func (bh *byH) pop() (h, i int) {
	h, i = bh.dh[0], bh.ix[0]
	n := len(bh.ix) - 1
	if n > 0 {
		bh.dh[n], bh.dh[0] = bh.dh[0], bh.dh[n]
		bh.ix[n], bh.ix[0] = bh.ix[0], bh.ix[n]
		bh.down(0, n)
	}
	bh.dh = bh.dh[:n]
	bh.ix = bh.ix[:n]
	return h, i
}

func (bh *byH) down(i0, n int) {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n &&
			bh.dh[j2] < bh.dh[j1] {
			j = j2 // = 2*i + 2  // right child
		}
		if !(bh.dh[j] < bh.dh[i]) {
			break
		}
		bh.dh[i], bh.dh[j] = bh.dh[j], bh.dh[i]
		bh.ix[i], bh.ix[j] = bh.ix[j], bh.ix[i]
		i = j
	}
}
