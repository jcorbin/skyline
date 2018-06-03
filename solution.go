package main

import (
	"container/heap"
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
}

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points. The returned point slice is
// only valid until the next call to Solve.
func (sol *Solver) Solve(data []internal.Building) ([]image.Point, error) {
	if len(data) == 0 {
		return nil, nil
	}
	bld := builder{
		res: make([]image.Point, 0, 1+len(data)*4),
	}
	pb := pending{
		bs: make([]internal.Building, 0, len(data)),
	}
	sort.Slice(data, func(i, j int) bool { return data[i].Sides[0] < data[j].Sides[0] })
	for _, b := range data {
		bld = bld.openBuilding(b, &pb)
	}
	bld = bld.closeOut(&pb)
	return bld.res, nil
}

type pending struct{ bs []internal.Building }

func (pb pending) Len() int            { return len(pb.bs) }
func (pb pending) Less(i, j int) bool  { return pb.bs[i].Sides[1] < pb.bs[j].Sides[1] }
func (pb pending) Swap(i, j int)       { pb.bs[i], pb.bs[j] = pb.bs[j], pb.bs[i] }
func (pb *pending) Push(x interface{}) { pb.bs = append(pb.bs, x.(internal.Building)) }
func (pb *pending) Pop() interface{} {
	i := len(pb.bs) - 1
	b := pb.bs[i]
	pb.bs = pb.bs[:i]
	return b
}

type builder struct {
	cur image.Point
	res []image.Point
}

func (bld builder) openBuilding(b internal.Building, pb *pending) builder {
	bld = bld.closePast(b.Sides[0], pb)
	if y := b.Height; y > bld.cur.Y {
		bld = bld.stepTo(b.Sides[0], y)
	}
	heap.Push(pb, b)
	return bld
}

func (bld builder) closePast(x int, pb *pending) builder {
	for pb.Len() > 0 && pb.bs[0].Sides[1] <= x {
		b := pb.bs[0]
		heap.Pop(pb)
		bld = bld.closeBuilding(b, pb)
	}
	return bld
}

func (bld builder) closeOut(pb *pending) builder {
	for pb.Len() > 0 {
		b := pb.bs[0]
		heap.Pop(pb)
		bld = bld.closeBuilding(b, pb)
	}
	return bld
}

func (bld builder) closeBuilding(b internal.Building, pb *pending) builder {
	if remHeight := maxHeightIn(pb.bs); remHeight < bld.cur.Y {
		bld = bld.stepTo(b.Sides[1], remHeight)
	}
	return bld
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

func (bld builder) stepTo(x, y int) builder {
	if x > bld.cur.X {
		bld = bld.tox(x)
	}
	bld = bld.toy(y)
	return bld
}

func (bld builder) tox(x int) builder { return bld.to(x, bld.cur.Y) }
func (bld builder) toy(y int) builder { return bld.to(bld.cur.X, y) }
func (bld builder) to(x, y int) builder {
	bld.cur.X = x
	bld.cur.Y = y
	bld.res = append(bld.res, bld.cur)
	return bld
}
