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
	o1  []int
	x1  []int
	x2  []int
	h   []int
	op  []pending
	res []image.Point
}

type pending struct{ x2, h, rh int }

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points. The returned point slice is
// only valid until the next call to Solve.
func (sol *Solver) Solve(data []internal.Building) ([]image.Point, error) {
	if len(data) == 0 {
		return nil, nil
	}
	o1, x1, x2, h := sol.load(data)
	return sol.run(o1, x1, x2, h)
}

func (sol *Solver) run(o1, x1, x2, h []int) ([]image.Point, error) {
	op := sol.op
	res := sol.res
	cur := image.ZP
	for o1i := 0; o1i < len(o1); o1i++ {
		// NOTE test probably doesn't catch edge case where several co-incident
		// opens cause redundant co-linear points
		i := o1[o1i]
		op, cur, res = closePast(i, x1, x2, op, cur, res)
		op, cur, res = open(i, x1, x2, h, op, cur, res)
	}
	op, cur, res = flush(x2, op, cur, res)
	return res, nil
}

func (sol *Solver) load(data []internal.Building) (_, _, _, _ []int) {
	sol.alloc(len(data))
	o1 := sol.o1
	x1 := sol.x1
	x2 := sol.x2
	h := sol.h
	for i := range data {
		o1, x1 = addOrderedXPoint(o1, x1, i, data[i].Sides[0])
		x2 = append(x2, data[i].Sides[1])
		h = append(h, data[i].Height)
	}
	return o1, x1, x2, h
}

func (sol *Solver) alloc(n int) {
	if cap(sol.o1) < n {
		sol.o1 = make([]int, 0, n)
		sol.x1 = make([]int, 0, n)
		sol.x2 = make([]int, 0, n)
		sol.h = make([]int, 0, n)
		sol.op = make([]pending, 0, n)
		sol.res = make([]image.Point, 0, 4*n)
	} else {
		sol.o1 = sol.o1[:0]
		sol.x1 = sol.x1[:0]
		sol.x2 = sol.x2[:0]
		sol.h = sol.h[:0]
		sol.op = sol.op[:0]
		sol.res = sol.res[:0]
	}
}

func open(
	i int, x1, x2, h []int,
	op []pending,
	cur image.Point, res []image.Point,
) (
	_ []pending,
	_ image.Point, _ []image.Point,
) {
	if bh := h[i]; bh > cur.Y {
		cur, res = tox(cur, res, x1[i])
		cur, res = goy(cur, res, bh)
	}
	op = appendRH(i, x2, h, op)
	return op, cur, res
}

func appendRH(i int, x2, h []int, op []pending) []pending {
	xi := x2[i]
	hi := h[i]

	// binary search for op-index where xi goes
	opi, nop := findRH(op, xi)

	// add new data at the end
	op = append(op, pending{xi, hi, 0})
	mh := op[opi].rh

	if opi != nop {
		// fix position of new data
		if oh := op[opi].h; mh < oh {
			mh = oh
		}
		copy(op[opi+1:], op[opi:])
		op[opi] = pending{xi, hi, mh}
	}

	// re-compute remaining height
	for opi > 0 {
		if oh := op[opi].h; mh < oh {
			mh = oh
		}
		opi--
		if mh > op[opi].rh {
			op[opi].rh = mh
		} else {
			break
		}
	}

	return op
}

func findRH(op []pending, x int) (_, _ int) {
	opi, nop := 0, len(op)
	for j := nop; opi < j; {
		h := int(uint(opi+j) >> 1)
		if op[h].x2 <= x {
			opi = h + 1
		} else {
			j = h
		}
	}
	return opi, nop
}

func closePast(
	i int, x1, x2 []int,
	op []pending,
	cur image.Point, res []image.Point,
) (
	_ []pending,
	_ image.Point, _ []image.Point,
) {
	bx := x1[i]
	opi := 0
	for ; opi < len(op); opi++ {
		x := op[opi].x2
		if x >= bx {
			break
		}
		if ah := op[opi].rh; ah < cur.Y {
			cur, res = tox(cur, res, x)
			cur, res = goy(cur, res, ah)
		}
	}
	op = op[:copy(op, op[opi:])]
	return op, cur, res
}

func flush(
	x2 []int,
	op []pending,
	cur image.Point, res []image.Point,
) (
	_ []pending,
	_ image.Point, _ []image.Point,
) {
	opi := 0
	for ; opi < len(op); opi++ {
		if ah := op[opi].rh; ah < cur.Y {
			cur, res = tox(cur, res, op[opi].x2)
			cur, res = goy(cur, res, ah)
		}
	}
	return op[:0], cur, res
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
