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

	o1 := sol.o1
	x1 := sol.x1
	x2 := sol.x2
	h := sol.h
	op := sol.op
	rh := sol.rh
	res := sol.res

	if cap(o1) < len(data) {
		o1 = make([]int, 0, len(data))
		x1 = make([]int, 0, len(data))
		x2 = make([]int, 0, len(data))
		h = make([]int, 0, len(data))
		op = make([]int, 0, len(data))
		rh = make([]int, 0, len(data))
		res = make([]image.Point, 0, 4*len(data))
	} else {
		o1 = o1[:0]
		x1 = x1[:0]
		x2 = x2[:0]
		h = h[:0]
		op = op[:0]
		rh = rh[:0]
		res = res[:0]
	}

	for i := range data {
		o1, x1 = addOrderedXPoint(o1, x1, i, data[i].Sides[0])
		x2 = append(x2, data[i].Sides[1])
		h = append(h, data[i].Height)
	}

	sol.o1 = o1
	sol.x1 = x1
	sol.x2 = x2
	sol.h = h
	sol.op = op
	sol.rh = rh

	sol.dir = dirNone
	sol.cur = image.ZP
	sol.res = res

	for o1i := 0; o1i < len(o1); o1i++ {
		// NOTE test probably doesn't catch edge case where several co-incident
		// opens cause redundant co-linear points
		i := o1[o1i]
		op, rh = sol.closePast(i, x1, x2, op, rh)
		op, rh = sol.open(i, x1, x2, h, op, rh)
	}
	op, rh = sol.flush(x2, op, rh)

	return sol.res, nil
}

func (sol *Solver) open(
	i int, x1, x2, h []int,
	op, rh []int,
) (_, _ []int) {
	if bh := h[i]; bh > sol.cur.Y {
		sol.tox(x1[i])
		sol.goy(bh)
	}
	op, rh = appendRH(i, x2, h, op, rh)
	return op, rh
}

func appendRH(i int, x2, h, op, rh []int) (_, _ []int) {
	// binary search for op-index where x2[i] goes
	opi, nop := findRH(x2, op, x2[i])

	// add new data at the end
	op, rh = append(op, i), append(rh, 0)
	mh := rh[opi]

	if opi != nop {
		// fix position of new data
		if oh := h[op[opi]]; mh < oh {
			mh = oh
		}
		copy(op[opi+1:], op[opi:])
		copy(rh[opi+1:], rh[opi:])
		op[opi] = i
		rh[opi] = mh
	}

	// re-compute remaining height
	for opi > 0 {
		if oh := h[op[opi]]; mh < oh {
			mh = oh
		}
		opi--
		if mh > rh[opi] {
			rh[opi] = mh
		} else {
			break
		}
	}

	return op, rh
}

func findRH(x2, op []int, x int) (_, _ int) {
	opi, nop := 0, len(op)
	for j := nop; opi < j; {
		h := int(uint(opi+j) >> 1)
		if x2[op[h]] <= x {
			opi = h + 1
		} else {
			j = h
		}
	}
	return opi, nop
}

func (sol *Solver) closePast(
	i int, x1, x2 []int,
	op, rh []int,
) (_, _ []int) {
	bx := x1[i]
	opi := 0
	for ; opi < len(op); opi++ {
		j := op[opi]
		if x2[j] >= bx {
			break
		}
		if ah := rh[opi]; ah < sol.cur.Y {
			sol.tox(x2[j])
			sol.goy(ah)
		}
	}
	op = op[:copy(op, op[opi:])]
	rh = rh[:copy(rh, rh[opi:])]
	return op, rh
}

func (sol *Solver) flush(
	x2 []int,
	op, rh []int,
) (_, _ []int) {
	opi := 0
	for ; opi < len(op); opi++ {
		j := op[opi]
		if ah := rh[opi]; ah < sol.cur.Y {
			sol.tox(x2[j])
			sol.goy(ah)
		}
	}
	op = op[:0]
	rh = rh[:0]
	return op, rh
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
