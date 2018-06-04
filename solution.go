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
	pb  pending
	bld builder
}

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points. The returned point slice is
// only valid until the next call to Solve.
func (sol *Solver) Solve(data []internal.Building) ([]image.Point, error) {
	if len(data) == 0 {
		return nil, nil
	}
	if n := 4*len(data) + 1; n > cap(sol.bld.res) {
		sol.bld.res = make([]image.Point, n)
	}
	if n := len(data); n > cap(sol.pb.co) {
		sol.pb = pending{
			co: make(closeOrder, n),
			rh: make([]int, n),
		}
	}
	sol.bld.cur = image.ZP
	sol.bld.res = sol.bld.res[:0]
	sol.pb.co = sol.pb.co[:0]
	sol.pb.rh = sol.pb.rh[:0]
	sort.Sort(openOrder(data))
	for _, b := range data {
		sol.pb = sol.bld.openBuilding(b, sol.pb)
	}
	sol.pb = sol.bld.closeOut(sol.pb)
	return sol.bld.res, nil
}

type openOrder []internal.Building
type closeOrder []internal.Building

func (oo openOrder) Len() int           { return len(oo) }
func (oo openOrder) Less(i, j int) bool { return oo[i].Sides[0] < oo[j].Sides[0] }
func (oo openOrder) Swap(i, j int)      { oo[i], oo[j] = oo[j], oo[i] }

func (co closeOrder) Len() int           { return len(co) }
func (co closeOrder) Less(i, j int) bool { return co[i].Sides[1] < co[j].Sides[1] }
func (co closeOrder) Swap(i, j int)      { co[i], co[j] = co[j], co[i] }

type pending struct {
	co closeOrder
	sx int
	rh []int
}

func (pb pending) find(i int) bool {
	return pb.co[i].Sides[1] > pb.sx
}

func (pb pending) append(b internal.Building) pending {
	pb.sx = b.Sides[1]
	n := len(pb.co)
	i := sort.Search(n, pb.find)
	pb.co = append(pb.co, b)
	if i != n {
		copy(pb.co[i+1:], pb.co[i:])
		pb.co[i] = b
	}
	pb.rh = append(pb.rh, 0)
	if i != n {
		copy(pb.rh[i+1:], pb.rh[i:])
		pb.rh[i] = pb.rh[i+1]
		// if i--; i >= 0 { pb.rh[i+1] = pb.rh[i] }
	} else {
		i--
	}

	for i--; i >= 0 && b.Height > pb.rh[i]; i-- {
		pb.rh[i] = b.Height
	}

	// TODO this should be the same, but it's not
	h := 0
	j := len(pb.rh) - 1
	for ; j >= 0; j-- {
		pb.rh[j] = h
		if jh := pb.co[j].Height; h < jh {
			h = jh
		}
	}

	return pb
}

func (pb pending) anyPast(x int) bool { return len(pb.co) > 0 && pb.co[0].Sides[1] <= x }

type builder struct {
	cur image.Point
	res []image.Point
}

func (bld *builder) openBuilding(b internal.Building, pb pending) pending {
	x := b.Sides[0]
	if pb.anyPast(x) {
		pb = bld.closePast(x, pb)
	}
	if y := b.Height; y > bld.cur.Y {
		bld.stepTo(x, y)
	}
	return pb.append(b)
}

func (bld *builder) closePast(x int, pb pending) pending {
	pb.sx = x
	j := sort.Search(len(pb.co), pb.find)
	for i := 0; i < j; i++ {
		bld.closeBuilding(i, pb)
	}
	pb.co = pb.co[:copy(pb.co, pb.co[j:])]
	pb.rh = pb.rh[:copy(pb.rh, pb.rh[j:])]
	return pb
}

func (bld *builder) closeOut(pb pending) pending {
	for i := 0; i < len(pb.co); i++ {
		bld.closeBuilding(i, pb)
	}
	pb.co = pb.co[:0]
	pb.rh = pb.rh[:0]
	return pb
}

func (bld *builder) closeBuilding(i int, pb pending) {
	if remHeight := pb.rh[i]; remHeight < bld.cur.Y {
		bld.stepTo(pb.co[i].Sides[1], remHeight)
	}
}

func (bld *builder) stepTo(x, y int) {
	if x > bld.cur.X {
		bld.tox(x)
	}
	bld.toy(y)
}

func (bld *builder) tox(x int) { bld.to(x, bld.cur.Y) }
func (bld *builder) toy(y int) { bld.to(bld.cur.X, y) }
func (bld *builder) to(x, y int) {
	bld.cur.X = x
	bld.cur.Y = y
	bld.res = append(bld.res, bld.cur)
}
