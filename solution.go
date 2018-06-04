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
		sol.pb = pending{co: make(closeOrder, n)}
	}
	sol.bld.cur = image.ZP
	sol.bld.res = sol.bld.res[:0]
	sol.pb.co = sol.pb.co[:0]
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
	mx int
}

func (pb pending) append(b internal.Building) pending {
	pb.co = append(pb.co, b)
	if pb.mx == 0 || b.Sides[1] < pb.mx {
		pb.mx = b.Sides[1]
	}
	return pb
}

func (pb pending) anyPast(x int) bool { return pb.mx <= x }

func (pb pending) heapify() {
	n := len(pb.co)
	for i := n/2 - 1; i >= 0; i-- {
		pb.down(i, n)
	}
}

func (pb pending) pop() (internal.Building, pending) {
	if i := len(pb.co) - 1; i > 0 {
		pb.co.Swap(0, i)
		pb.down(0, i)
		b := pb.co[i]
		pb.mx = pb.co[0].Sides[1]
		pb.co = pb.co[:i]
		return b, pb
	}
	b := pb.co[0]
	return b, pending{co: pb.co[:0]}
}

func (pb pending) down(i0, n int) {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && pb.co.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !pb.co.Less(j, i) {
			break
		}
		pb.co.Swap(i, j)
		i = j
	}
}

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
	pb.heapify()
	for len(pb.co) > 0 && pb.co[0].Sides[1] <= x {
		var b internal.Building
		b, pb = pb.pop()
		bld.closeBuilding(b, pb)
	}
	return pb
}

func (bld *builder) closeOut(pb pending) pending {
	pb.heapify()
	for len(pb.co) > 0 {
		var b internal.Building
		b, pb = pb.pop()
		bld.closeBuilding(b, pb)
	}
	return pb
}

func (bld *builder) closeBuilding(b internal.Building, pb pending) {
	if remHeight := maxHeightIn(pb.co); remHeight < bld.cur.Y {
		bld.stepTo(b.Sides[1], remHeight)
	}
}

// maxHeightIn computes the maximum height in a slice of buildings; it is used
// in context to compute the "remaining pending height", and as such is doing
// so in a wildly inefficient manner:
// - since the caller will call this utility for each buildings[i] on
//   buildings[i+1:], we're continually re-computing suffix-max-heights
// - it'd be better to iterate buildings once in reverse order, collecting
//   cum-max-heights
func maxHeightIn(buildings []internal.Building) (h int) {
	for j := 0; j < len(buildings); j++ {
		if bh := buildings[j].Height; bh > h {
			h = bh
		}
	}
	return h
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
