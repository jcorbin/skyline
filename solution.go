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

	res = mergeWall(res, li, x0, 0, h)

	if ri == len(res) {
		return append(res,
			image.Pt(x1, h), image.Pt(x1, 0))
	}

	res = mergeWall(res, ri, x1, h, 0)

	return res
}

// mergeWall at into an x-sorted result array.
// - i is the index where the new x value lands; there's either one or more
//   existing points with equal X value, or the index is of the first point
//   beyond x.
// - x is the X value that any newly added points should have
// - y0 and y1 describe the wall's vertical span; their order matters (i.e. it
//   is the case that either y0 < y1 or y1 < y0), and any points added should
//   maintain their y ordering.  NOTE order of y0 and y1 matters
func mergeWall(res []image.Point, i, x, y0, y1 int) []image.Point {
	// resolve new wall parity
	rising := y1 >= y0
	if !rising {
		y0, y1 = y1, y0
	}

	// resolve prior wall parity
	j := i     // index of the wall bottom
	k := j + 1 // index of wall top
	i--        // index of prior roof start
	oldRising := res[k].Y >= res[j].Y
	if !oldRising {
		j, k = k, j
	}

	if res[j].X == x {
		// co-linear prior wall
		if y1 > res[k].Y {
			res[k].Y, res[k+1].Y = y1, y1 // raise the roof
		}
		return res
	}

	if res[j].Y > y1 {
		// wall wholly occluded
		return res
	}

	// TODO describe; refactor / reconsider control flow once we get through
	// case enumeration

	if oldRising {

		if rising {
			if y1 < res[k].Y {
				res = append(res, image.ZP, image.ZP)
				copy(res[k+2:], res[k:])
				/*
				*  Before	    After
				*		  k	         k+2
				*		  |	         |
				*		  |	      k--k+1
				*		  |	      |
				* __i_____j  __i__j
				 */
				res[j].X = x
				res[k].X = x
				res[k].Y = y1
				res[k+1].X = res[k+2].X
				res[k+1].Y = y1
			} else {
				/*
				*  Before	    After
				*	         	    k--k+1
				*	         	    |
				*	    k-k+1	    |
				*	    |    	    |
				* __i___j      __i__j
				 */
				res[j].X = x
				res[k].X = x
				res[k].Y = y1
				res[k+1].Y = y1
			}
		} else {

			if y1 < res[k].Y {
				/*
				*  Before	    After
				* __i_____j  TODO
				*		  |	 TODO
				*		  |	 TODO
				*		  |	 TODO
				*		  k	 TODO
				 */

			} else {
				// FIXME
			}
			panic("unimplemented")

		}

	}

	return res
}
