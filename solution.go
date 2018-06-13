package main

import (
	"image"
	"log"

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

		found := false // TODO for the case of full prefix when binary searching

		for j, pb := range prior {
			log.Printf("prior[%v] = %v", j, pb)
			if pb.Sides[1] < x1 {
				// TODO binary search for this, and copy() the prefix
				next = append(next, pb)
				continue
			}

			if pb.Sides[0] > x2 {
				// TODO binary search for this, and copy() the suffix
				next = append(next, pb)
				continue
			}

			if h > pb.Height {
				found = true
				if pb.Sides[0] < x1 {
					next = append(next, internal.Bldg(pb.Sides[0], x1, pb.Height))
				}
				next = append(next, internal.Bldg(x1, x2, h))
				if pb.Sides[1] > x2 {
					next = append(next, internal.Bldg(x2, pb.Sides[1], pb.Height))
				}
			}
		}
		if !found {
			next = append(next, b)
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
