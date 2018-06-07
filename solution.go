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
	o1 []int
	x1 []int
	x2 []int
	h  []int
	op []int
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
		sol.o1, sol.x1 = addOrderedXPoint(sol.o1, sol.x1, i, data[i].Sides[0])
		sol.x2 = append(sol.x2, data[i].Sides[1])
		sol.h = append(sol.h, data[i].Height)
	}

	sol.dir = dirNone
	sol.cur = image.ZP
	for o1i := 0; o1i < len(sol.o1); o1i++ {
		// NOTE test probably doesn't catch edge case where several co-incident
		// opens cause redundant co-linear points
		i := sol.o1[o1i]
		sol.closePast(i)
		sol.open(i)
	}
	sol.flush()

	return sol.res, nil
}

func (sol *Solver) alloc(n int) {
	if cap(sol.o1) < n {
		sol.o1 = make([]int, 0, n)
		sol.x1 = make([]int, 0, n)
		sol.x2 = make([]int, 0, n)
		sol.h = make([]int, 0, n)
		sol.op = make([]int, 0, n)
		sol.rh = make([]int, 0, n)
		sol.res = make([]image.Point, 0, 4*n)
	} else {
		sol.o1 = sol.o1[:0]
		sol.x1 = sol.x1[:0]
		sol.x2 = sol.x2[:0]
		sol.h = sol.h[:0]
		sol.op = sol.op[:0]
		sol.rh = sol.rh[:0]
		sol.res = sol.res[:0]
	}
}

func (sol *Solver) open(i int) {
	if bh := sol.h[i]; bh > sol.cur.Y {
		sol.tox(sol.x1[i])
		sol.goy(bh)
	}
	sol.appendRH(i)
}

func (sol *Solver) appendRH(i int) {
	// binary search for op-index where x2[i] goes
	opi, nop := sol.findRH(sol.x2[i])

	// add new data at the end
	sol.op, sol.rh = append(sol.op, i), append(sol.rh, 0)
	mh := sol.rh[opi]

	if opi != nop {
		// fix position of new data
		if oh := sol.h[sol.op[opi]]; mh < oh {
			mh = oh
		}
		copy(sol.op[opi+1:], sol.op[opi:])
		copy(sol.rh[opi+1:], sol.rh[opi:])
		sol.op[opi] = i
		sol.rh[opi] = mh
	}

	// re-compute remaining height
	for opi > 0 {
		if oh := sol.h[sol.op[opi]]; mh < oh {
			mh = oh
		}
		opi--
		if mh > sol.rh[opi] {
			sol.rh[opi] = mh
		} else {
			break
		}
	}
}

func (sol *Solver) findRH(x int) (_, _ int) {
	opi, nop := 0, len(sol.op)
	for j := nop; opi < j; {
		h := int(uint(opi+j) >> 1)
		if sol.x2[sol.op[h]] <= x {
			opi = h + 1
		} else {
			j = h
		}
	}
	return opi, nop
}

func (sol *Solver) closePast(i int) {
	bx := sol.x1[i]
	opi := 0
	for ; opi < len(sol.op); opi++ {
		j := sol.op[opi]
		if sol.x2[j] >= bx {
			break
		}
		if ah := sol.rh[opi]; ah < sol.cur.Y {
			sol.tox(sol.x2[j])
			sol.goy(ah)
		}
	}
	sol.op = sol.op[:copy(sol.op, sol.op[opi:])]
	sol.rh = sol.rh[:copy(sol.rh, sol.rh[opi:])]
}

func (sol *Solver) flush() {
	opi := 0
	for ; opi < len(sol.op); opi++ {
		j := sol.op[opi]
		if ah := sol.rh[opi]; ah < sol.cur.Y {
			sol.tox(sol.x2[j])
			sol.goy(ah)
		}
	}
	sol.op = sol.op[:0]
	sol.rh = sol.rh[:0]
}

func addOrderedXPoint(os, xs []int, i, x int) (_, _ []int) {
	oi, on := findXPoint(os, xs, x)
	xs = append(xs, x)
	if oi == on {
		os = append(os, i)
	} else {
		os = append(os, i)
		copy(os[oi+1:], os[oi:])
		os[oi] = i
	}
	return os, xs
}

func findXPoint(os, xs []int, x int) (_, _ int) {
	oi := 0
	on := len(os)
	for j := on; oi < j; {
		h := int(uint(oi+j) >> 1)
		if xs[os[h]] <= x {
			oi = h + 1
		} else {
			j = h
		}
	}
	return oi, on
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
