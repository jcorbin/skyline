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
	sol.bld = makeBuilder(1 + len(data)*4)
	sol.pb = makePending(len(data))
	sort.Slice(data, func(i, j int) bool { return data[i].Sides[0] < data[j].Sides[0] })
	for _, b := range data {
		sol.pb = sol.bld.openBuilding(b, sol.pb)
	}
	sol.pb = sol.bld.closeOut(sol.pb)
	return sol.bld.res, nil
}

type pending []internal.Building

func makePending(cap int) pending     { return make(pending, 0, cap) }
func (pb pending) less(i, j int) bool { return pb[i].Sides[1] < pb[j].Sides[1] }
func (pb pending) swap(i, j int)      { pb[i], pb[j] = pb[j], pb[i] }

func (pb pending) anyPast(x int) bool {
	for i := range pb {
		if pb[i].Sides[1] <= x {
			return true
		}
	}
	return false
}

func (pb pending) heapify() {
	n := len(pb)
	for i := n/2 - 1; i >= 0; i-- {
		pb.down(i, n)
	}
}

func (pb pending) pop() (internal.Building, pending) {
	i := len(pb) - 1
	pb.swap(0, i)
	pb.down(0, i)
	b := pb[i]
	return b, pb[:i]
}

func (pb pending) down(i0, n int) {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && pb.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !pb.less(j, i) {
			break
		}
		pb.swap(i, j)
		i = j
	}
}

type builder struct {
	cur image.Point
	res []image.Point
}

func makeBuilder(cap int) builder {
	return builder{
		res: make([]image.Point, 0, cap),
	}
}

func (bld *builder) openBuilding(b internal.Building, pb pending) pending {
	x := b.Sides[0]
	if pb.anyPast(x) {
		pb = bld.closePast(x, pb)
	}
	if y := b.Height; y > bld.cur.Y {
		bld.stepTo(x, y)
	}
	return append(pb, b)
}

func (bld *builder) closePast(x int, pb pending) pending {
	pb.heapify()
	for len(pb) > 0 && pb[0].Sides[1] <= x {
		var b internal.Building
		b, pb = pb.pop()
		bld.closeBuilding(b, pb)
	}
	return pb
}

func (bld *builder) closeOut(pb pending) pending {
	pb.heapify()
	for len(pb) > 0 {
		var b internal.Building
		b, pb = pb.pop()
		bld.closeBuilding(b, pb)
	}
	return pb
}

func (bld *builder) closeBuilding(b internal.Building, pb pending) {
	if remHeight := maxHeightIn(pb); remHeight < bld.cur.Y {
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
