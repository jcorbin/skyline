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
	if len(data) == 0 {
		return nil, nil
	}
	bld := makeBuilder(1 + len(data)*4)
	pb := makePending(len(data))
	sort.Slice(data, func(i, j int) bool { return data[i].Sides[0] < data[j].Sides[0] })
	for _, b := range data {
		bld = bld.openBuilding(b, &pb)
	}
	bld = bld.closeOut(&pb)
	return bld.res, nil
}

type pending struct{ bs []internal.Building }

func makePending(cap int) pending      { return pending{make([]internal.Building, 0, cap)} }
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
func (pb pending) AnyPast(x int) bool {
	for i := range pb.bs {
		if pb.bs[i].Sides[1] <= x {
			return true
		}
	}
	return false
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

func (bld builder) openBuilding(b internal.Building, pb *pending) builder {
	x := b.Sides[0]
	if pb.AnyPast(x) {
		bld = bld.closePast(x, pb)
	}
	if y := b.Height; y > bld.cur.Y {
		bld = bld.stepTo(x, y)
	}
	pb.bs = append(pb.bs, b)
	return bld
}

func (bld builder) closePast(x int, pb *pending) builder {
	heap.Init(pb)
	for pb.Len() > 0 && pb.bs[0].Sides[1] <= x {
		b := pb.bs[0]
		heap.Pop(pb)
		bld = bld.closeBuilding(b, pb)
	}
	return bld
}

func (bld builder) closeOut(pb *pending) builder {
	heap.Init(pb)
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
