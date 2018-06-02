package main

import (
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
	pending := make([]internal.Building, 0, len(data))
	sort.Slice(data, func(i, j int) bool { return data[i].Sides[0] < data[j].Sides[0] })
	for _, b := range data {
		bld, pending = bld.openBuilding(b, pending)
	}
	bld, pending = bld.closePast(0, pending)
	return bld.res, nil
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

func (bld builder) openBuilding(b internal.Building, pending []internal.Building) (builder, []internal.Building) {
	bld, pending = bld.closePast(b.Sides[0], pending)
	if y := b.Height; y > bld.cur.Y {
		bld = bld.stepTo(b.Sides[0], y)
	}
	pending = append(pending, b)
	sort.Slice(pending, func(i, j int) bool { return pending[i].Sides[1] < pending[j].Sides[1] })
	return bld, pending
}

func (bld builder) closePast(x int, pending []internal.Building) (builder, []internal.Building) {
	i := 0
	for ; i < len(pending) && (x <= 0 || pending[i].Sides[1] <= x); i++ {
		if remHeight := maxHeightIn(pending[i+1:]); remHeight < bld.cur.Y {
			bld = bld.stepTo(pending[i].Sides[1], remHeight)
		}
	}
	pending = pending[:copy(pending, pending[i:])]
	return bld, pending
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
