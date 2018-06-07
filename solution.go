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
	x1 []int
	x2 []int
	h  []int

	ix     []int
	h1, h2 xHeap
	// rh     []int

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
		sol.ix = append(sol.ix, i)
	}

	sol.h1 = xHeap{sol.ix, sol.x1}     // nil, nil
	sol.h2 = xHeap{sol.ix[:0], sol.x2} // sol.rh, sol.h
	sol.h1.heapify()

	sol.dir = dirNone
	sol.cur = image.ZP
	for len(sol.h1.ix) > 0 {
		bx := sol.h1.minx()
		for len(sol.h2.ix) > 0 {
			ax := sol.h2.minx()
			if ax >= bx {
				break
			}

			ah := 0
			for _, i := range sol.h2.ix[1:] {
				if h := sol.h[i]; h > ah {
					ah = h
				}
			}
			// ah := sol.h2.maxh()

			if ah < sol.cur.Y {
				sol.tox(ax)
				sol.goy(ah)
			}
			sol.h2.shift()
		}

		i := sol.h1.ix[0]
		if bh := sol.h[i]; bh > sol.cur.Y {
			sol.tox(bx)
			sol.goy(bh)
		}

		sol.h1.shift()
		sol.h2.subsume()
	}
	// log.Printf("h2: %+v", sol.h2)

	for len(sol.h2.ix) > 0 {

		ah := 0
		for _, i := range sol.h2.ix[1:] {
			if h := sol.h[i]; h > ah {
				ah = h
			}
		}
		// ah := sol.h2.maxh()

		if ah < sol.cur.Y {
			sol.tox(sol.h2.minx())
			sol.goy(ah)
		}
		sol.h2.shift()
	}

	return sol.res, nil
}

func (sol *Solver) alloc(n int) {
	if cap(sol.x1) < n {
		sol.x1 = make([]int, 0, n)
		sol.x2 = make([]int, 0, n)
		sol.h = make([]int, 0, n)
		sol.ix = make([]int, 0, n)
		// sol.rh = make([]int, 0, n)
		sol.res = make([]image.Point, 0, 4*n)
	} else {
		sol.x1 = sol.x1[:0]
		sol.x2 = sol.x2[:0]
		sol.h = sol.h[:0]
		sol.ix = sol.ix[:0]
		// sol.rh = sol.rh[:0]
		sol.res = sol.res[:0]
	}
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

type xHeap struct {
	ix []int
	xs []int
	// rh []int
	// h  []int
}

func (xh xHeap) minx() int { return xh.xs[xh.ix[0]] }

// func (xh xHeap) maxh() int { return xh.rh[0] }

func (xh *xHeap) shift() {
	xh.ix = xh.ix[1:]
	xh.heapify() // TODO be-headed; merge two valid sub-heaps
}

func (xh *xHeap) subsume() {
	if i := len(xh.ix); i < cap(xh.ix) {
		xh.ix = xh.ix[:i+1]
		// xh.rh = append(xh.rh, 0)
		xh.up(i)
	}
}

// TODO func (xh xHeap) merge() when composed of two valid heaps (e.g. after be-heading)

func (xh xHeap) heapify() {
	n := len(xh.ix)
	for i := n/2 - 1; i >= 0; i-- {
		xh.down(i, n)
	}
}

func (xh xHeap) less(i, j int) bool { return xh.xs[xh.ix[i]] < xh.xs[xh.ix[j]] }
func (xh xHeap) swap(i, j int)      { xh.ix[i], xh.ix[j] = xh.ix[j], xh.ix[i] }

func (xh xHeap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !xh.less(j, i) {
			break
		}
		xh.swap(i, j)

		// if xh.rh != nil {
		// 	// TODO can we do less?
		// 	xh.fixrh(i)
		// 	xh.fixrh(j)
		// 	// TODO towards less`
		// 	// h := xh.rh[j]
		// 	// if h == 0 {
		// 	// 	h = xh.h[xh.ix[j]]
		// 	// }
		// 	// if h > xh.rh[i] {
		// 	// 	xh.rh[i] = h
		// 	// }
		// 	// xh.fixrh(i)
		// }

		j = i
	}
}

func (xh xHeap) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && xh.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !xh.less(j, i) {
			break
		}
		xh.swap(i, j)

		// if xh.rh != nil {
		// 	// TODO can we do less?
		// 	xh.fixrh(i)
		// 	xh.fixrh(j)
		// }

		i = j
	}
	return i > i0
}

// func (xh xHeap) fixrh(i int) {
// 	h := 0
// 	n := len(xh.ix)
// 	if j1 := 2*i + 1; j1 < n {
// 		h = max3(h, xh.h[xh.ix[j1]], xh.rh[j1])
// 		if j2 := j1 + 1; j2 < n {
// 			h = max3(h, xh.h[xh.ix[j2]], xh.rh[j2])
// 		}
// 	}
// 	xh.rh[i] = h
// }

// func max3(a, b, c int) int {
// 	if b > a {
// 		a = b
// 	}
// 	if c > a {
// 		a = c
// 	}
// 	return a
// }
