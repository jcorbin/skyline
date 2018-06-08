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
	open []internal.Building

	cur image.Point
	res []image.Point
}

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points. The returned point slice is
// only valid until the next call to Solve.
func (sol *Solver) Solve(data []internal.Building) ([]image.Point, error) {
	sol.alloc(len(data))

	sort.Slice(data, func(i, j int) bool {
		return data[i].Sides[0] < data[j].Sides[0]
	})

	for _, b := range data {
		if len(sol.open) > 0 && sol.open[len(sol.open)-1].Sides[1] < b.Sides[0] {
			sol.closeOpen()
		}
		sol.openBuilding(b)
	}
	sol.closeOpen()

	return sol.res, nil
}

func (sol *Solver) openBuilding(b internal.Building) {
	log.Printf("open %v", b)
	x2 := b.Sides[1]
	i := sort.Search(len(sol.open), func(i int) bool {
		return sol.open[i].Sides[1] > x2
	})
	sol.open = append(sol.open, b)
	if i != len(sol.open) {
		copy(sol.open[i+1:], sol.open[i:])
		sol.open[i] = b
	}
	if b.Height > sol.cur.Y {
		sol.tox(b.Sides[0])
		sol.goy(b.Height)
	}
}

func (sol *Solver) closeOpen() {
	rh := calcRemHeight(sol.open)
	log.Printf("close: %v", sol.open)
	log.Printf("rh: %v", rh)
	for i, b := range sol.open {
		log.Printf("- [%v] %v", i, b)
		if h := rh[i]; h != sol.cur.Y {
			sol.tox(b.Sides[1])
			sol.goy(h)
		}
	}
}

func (sol *Solver) alloc(n int) {
	if m := 4 * n; m < cap(sol.res) {
		sol.res = make([]image.Point, 0, m)
	} else {
		sol.res = sol.res[:0]
	}
}

func (sol *Solver) tox(x int) {
	if x != sol.cur.X {
		sol.gox(x)
	}
}

func (sol *Solver) gox(x int) {
	sol.cur.X = x
	if runningy(sol.res) {
		sol.res[len(sol.res)-1].X = x
	} else {
		sol.res = append(sol.res, sol.cur)
	}
}

func (sol *Solver) goy(y int) {
	sol.cur.Y = y
	if runningx(sol.res) {
		sol.res[len(sol.res)-1].Y = y
	} else {
		sol.res = append(sol.res, sol.cur)
	}
}

func runningx(pts []image.Point) bool {
	i := len(pts) - 2
	j := i + 1
	if i >= 0 {
		if pts[i].X == pts[j].X {
			return true
		}
	}
	return false
}

func runningy(pts []image.Point) bool {
	i := len(pts) - 2
	j := i + 1
	if i >= 0 {
		if pts[i].Y == pts[j].Y {
			return true
		}
	}
	return false
}
