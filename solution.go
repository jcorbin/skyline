package main

import (
	"image"
	"log"
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
	cur image.Point
	res []image.Point
}

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points. The returned point slice is
// only valid until the next call to Solve.
func (sol *Solver) Solve(data []internal.Building) ([]image.Point, error) {
	if len(data) == 0 {
		return nil, nil
	}

	est := len(data)
	next := make([]internal.Building, 0, est)
	prior := append(make([]internal.Building, 0, est), data[0])
	sol.res = make([]image.Point, 0, 4*len(data))

	for i := 1; i < len(data); i++ {
		b := data[i]

		log.Printf("add %v", b)
		x1, x2, h := b.Sides[0], b.Sides[1], b.Height

		ix1 := sort.Search(len(prior), func(i int) bool { return x1 >= prior[i].Sides[0] })
		ix2 := sort.Search(len(prior), func(i int) bool { return x2 <= prior[i].Sides[1] })

		// new building just needs to be appended to end
		if ix1 == len(prior) {
			next = append(next, prior...)
			next = append(next, b)
			log.Printf("APP %v", prior)
			continue
		}

		// prefix before new building
		if ix1 > 0 {
			next = append(next, prior[:ix1]...)
		}

		// new building overlaps any buildings in prior[ix1:ix2]

		// (chunks of) prior building(s) will make it into next if:
		// - the first one can have a hanging left section
		// - anything in the middle can poke out above
		// - the last one can have a hanging right section

		// - the first one can have a hanging left section
		if ix1 < len(prior) {
			if px := prior[ix1].Sides[0]; px < x1 {
				next = append(next, internal.Bldg(px, x1, prior[ix1].Height))
				ix1++ // XXX what about right poking out too?
			}
		}
		// prior[ix1].Sides[0] >= x1

		// - anything in the middle can poke out above
		// for ix1 <= ix2 && ix1 < len(prior) { }
		// FIXME

		// - the last one can have a hanging right section
		if ix2 < len(prior) {
			if px := prior[ix2].Sides[1]; px > x2 {
				// next = append(next, internal.Bldg(px, x1, prior[ix1].Height)) XXX
				ix2++
			}
		}
		// FIXME

		// if ix2 == len(prior) { } XXX ?

		// any remaining chunk (maybe the entire new building if no overlap)
		if x1 < x2 {
			next = append(next, internal.Bldg(x1, x2, h))
		}

		// suffix after new building
		if ixsuf := ix2 + 1; ixsuf < len(prior) {
			next = append(next, prior[ixsuf:]...)
		}

		prior, next = next, prior[:0]
		log.Printf("added %v", prior)
	}

	log.Printf("proc %v", prior)

	for _, b := range prior {
		sol.tox(b.Sides[0])
		sol.toy(b.Height)
		sol.gox(b.Sides[1])
		sol.goy(0)
	}

	return sol.res, nil
}

func (sol *Solver) tox(x int) {
	if sol.cur.X != x {
		sol.gox(x)
	}
}
func (sol *Solver) toy(y int) {
	if sol.cur.Y != y {
		sol.goy(y)
	}
}

func (sol *Solver) gox(x int) {
	sol.cur.X = x
	if j := len(sol.res) - 1; j > 0 {
		if i := j - 1; sol.res[i].Y == sol.res[j].Y {
			// log.Printf("cox(%v) %v", x, sol.res)
			sol.res[j].X = x
			return
		}
	}
	sol.res = append(sol.res, sol.cur)
	// log.Printf("gox(%v) %v", x, sol.res)
}

func (sol *Solver) goy(y int) {
	sol.cur.Y = y
	if j := len(sol.res) - 1; j > 0 {
		if i := j - 1; sol.res[i].X == sol.res[j].X {
			// log.Printf("coy(%v) %v", y, sol.res)
			sol.res[j].Y = y
			return
		}
	}
	sol.res = append(sol.res, sol.cur)
	// log.Printf("goy(%v) %v", y, sol.res)
}
