package main

import (
	"image"
	"log"

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

	o  []int
	c  []int
	rh []int

	dir bldDir
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

	sol.alloc(len(data))

	for i := range data {
		sol.x1 = append(sol.x1, data[i].Sides[0])
		sol.x2 = append(sol.x2, data[i].Sides[1])
		sol.h = append(sol.h, data[i].Height)
		sol.o = append(sol.o, i)
	}

	heapify(sol.o, sol.x1)

	sol.dir = dirNone
	sol.cur = image.ZP
	for len(sol.o) > 0 {
		i := sol.o[0]

		// closePast
		bx := sol.x1[i]
		for len(sol.c) > 0 {
			j := sol.c[0]
			if sol.x2[j] >= bx {
				break
			}
			sol.c = heappop(sol.c, sol.x2)
			sol.recomputerh(0)
			if ah := sol.rh[j]; ah < sol.cur.Y {
				sol.tox(sol.x2[j])
				sol.goy(ah)
			}
		}

		// open
		if bh := sol.h[i]; bh > sol.cur.Y {
			sol.tox(sol.x1[i])
			sol.goy(bh)
		}
		sol.o, sol.c = heapshift(sol.o, sol.x1, sol.c, sol.x2)
		sol.recomputerh(0)
	}

	// flush
	for len(sol.c) > 0 {
		j := sol.c[0]
		log.Printf("flush [%v] %v", j, data[j])
		sol.c = heappop(sol.c, sol.x2)
		sol.recomputerh(0)
		if ah := sol.rh[j]; ah <= sol.cur.Y {
			sol.tox(sol.x2[j])
			sol.goy(ah)
		}
	}

	return sol.res, nil
}

func (sol *Solver) alloc(n int) {
	if cap(sol.x1) < n {
		sol.x1 = make([]int, 0, n)
		sol.x2 = make([]int, 0, n)
		sol.h = make([]int, 0, n)
		sol.o = make([]int, 0, n)
		sol.rh = make([]int, n)
		sol.res = make([]image.Point, 0, 4*n)
	} else {
		sol.x1 = sol.x1[:0]
		sol.x2 = sol.x2[:0]
		sol.h = sol.h[:0]
		sol.o = sol.o[:0]
		sol.rh = sol.rh[:n]
		sol.res = sol.res[:0]
	}
	sol.c = sol.o[len(sol.o):]
}

func (sol *Solver) recomputerh(ci int) {
	rh := 0
	for _, child := range []int{2*ci + 1, 2*ci + 2} {
		if child >= len(sol.c) {
			break
		}
		sol.recomputerh(child)
		id := sol.c[child]
		rh = max3(rh, sol.h[id], sol.rh[id])
	}
}

func max3(a, b, c int) int {
	if a < b {
		a = b
	}
	if a < c {
		a = c
	}
	return a
}

type bldDir uint8

const (
	dirNone bldDir = iota
	dirVert
	dirHoriz
)

func (sol *Solver) goy(y int) {
	sol.cur.Y = y
	if sol.dir == dirVert {
		sol.res[len(sol.res)-1].Y = y
	} else {
		sol.res = append(sol.res, sol.cur)
	}
	sol.dir = dirVert
}

func (sol *Solver) tox(x int) {
	if x != sol.cur.X {
		sol.cur.X = x
		if sol.dir == dirHoriz {
			sol.res[len(sol.res)-1].X = x
		} else {
			sol.res = append(sol.res, sol.cur)
		}
		sol.dir = dirHoriz
	}
}

func heapify(ix, xs []int) {
	n := len(ix)
	for i := n/2 - 1; i >= 0; i-- {
		heapdown(ix, xs, i, n)
	}
}

func heapdown(ix, xs []int, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && xs[ix[j2]] < xs[ix[j1]] {
			j = j2 // = 2*i + 2  // right child
		}
		if xs[ix[j]] >= xs[ix[i]] {
			break
		}
		ix[i], ix[j] = ix[j], ix[i]
		i = j
	}
	return i > i0
}

func heappop(ix, xs []int) []int {
	i := len(ix) - 1
	if i > 0 {
		ix[0], ix[i] = ix[i], ix[0]
		heapdown(ix, xs, 0, i)
	}
	return ix[:i]
}

func heapshift(ix1, x1s, ix2, x2s []int) (_, _ []int) {
	i := len(ix1) - 1
	if i > 0 {
		ix1[0], ix1[i] = ix1[i], ix1[0]
		heapdown(ix1, x1s, 0, i)
	}
	ix1, ix2 = ix1[:i], ix1[i:i+len(ix2)+1]
	heapdown(ix2, x2s, 0, len(ix2))
	return ix1, ix2
}
