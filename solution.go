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
	x2  []int
	h   []int
	op  []int
	rh  []int
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
	sol.res = res

	cur := image.ZP
	for o1i := 0; o1i < len(o1); o1i++ {
		// NOTE test probably doesn't catch edge case where several co-incident
		// opens cause redundant co-linear points
		i := o1[o1i]
		op, rh, cur, res = closePast(i, x1, x2, op, rh, cur, res)
		op, rh, cur, res = open(i, x1, x2, h, op, rh, cur, res)
	}
	op, rh, cur, res = flush(x2, op, rh, cur, res)

	return res, nil
}

func open(
	i int, x1, x2, h []int,
	op, rh []int,
	cur image.Point, res []image.Point,
) (
	_, _ []int,
	_ image.Point, _ []image.Point,
) {
	if bh := h[i]; bh > cur.Y {
		cur, res = tox(cur, res, x1[i])
		cur, res = goy(cur, res, bh)
	}
	op, rh = appendRH(i, x2, h, op, rh)
	return op, rh, cur, res
}

func appendRH(i int, x2, h, op, rh []int) (_, _ []int) {
	nop := len(op)
	opi := sort.Search(nop, func(opi int) bool { return x2[op[opi]] > x2[i] })
	op, rh = append(op, i), append(rh, 0)
	mh := rh[opi]

	if opi != nop {
		if oh := h[op[opi]]; mh < oh {
			mh = oh
		}
		copy(op[opi+1:], op[opi:])
		copy(rh[opi+1:], rh[opi:])
		op[opi] = i
		rh[opi] = mh
	}

	for opi > 0 {
		if oh := h[op[opi]]; mh < oh {
			mh = oh
		}
		opi--
		rh[opi] = mh
		// TODO should be able to halt early
		// if rh[opi] < bh { rh[opi] = bh }
	}

	return op, rh
}

func closePast(
	i int, x1, x2 []int,
	op, rh []int,
	cur image.Point, res []image.Point,
) (
	_, _ []int,
	_ image.Point, _ []image.Point,
) {
	bx := x1[i]
	opi := 0
	for ; opi < len(op); opi++ {
		j := op[opi]
		if x2[j] >= bx {
			break
		}
		if ah := rh[opi]; ah < cur.Y {
			cur, res = tox(cur, res, x2[j])
			cur, res = goy(cur, res, ah)
		}
	}
	op = op[:copy(op, op[opi:])]
	rh = rh[:copy(rh, rh[opi:])]
	return op, rh, cur, res
}

func flush(
	x2 []int,
	op, rh []int,
	cur image.Point, res []image.Point,
) (
	_, _ []int,
	_ image.Point, _ []image.Point,
) {
	opi := 0
	for ; opi < len(op); opi++ {
		j := op[opi]
		if ah := rh[opi]; ah < cur.Y {
			cur, res = tox(cur, res, x2[j])
			cur, res = goy(cur, res, ah)
		}
	}
	op = op[:0]
	rh = rh[:0]
	return op, rh, cur, res
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
