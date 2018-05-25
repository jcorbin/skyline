package main

import (
	"image"
	"sort"

	"github.com/jcorbin/skyline/internal"
)

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points.
func Solve(data []internal.Building) ([]image.Point, error) {
	res := make([]image.Point, 1, 1+len(data)*4)
	res[0] = image.ZP
	for _, b := range data {
		res = mergeBuilding(res, b.Sides[0], b.Sides[1], b.Height)
	}
	return res, nil
}

// mergeBuilding into an x-sorted result array.
func mergeBuilding(res []image.Point, x0, x1, h int) []image.Point {
	var (
		li = sort.Search(len(res), func(i int) bool { return res[i].X >= x0 })
		ri = sort.Search(len(res), func(i int) bool { return res[i].X >= x1 })
	)

	if li == len(res) {
		return append(res,
			image.Pt(x0, 0), image.Pt(x0, h),
			image.Pt(x1, h), image.Pt(x1, 0))
	}

	if ri == len(res) {
		panic("unimplemented extend last")
		// return extendLast(res, x0, x1, h)
	}

	/*
	 *  +---+
	 *  | h |
	 * -x0 x1-
	 */

	/*
	 *  +--ri
	 *  |   |
	 * -li  +-
	 */

	/*
	 * -li  +-
	 *   |  |
	 *   +-ri
	 */

	/*
	 *  +--+
	 *  |  |
	 * -li | +-
	 *     | |
	 *     +-ri
	 */

	/*
	 *     +-ri
	 *     |  |
	 * -li |  +-
	 *   | |
	 *   +-+
	 */

	return res
}
