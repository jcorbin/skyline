package main

import (
	"fmt"
	"image"
	"log"
	"sort"

	"github.com/jcorbin/skyline/internal"
)

// Solve receives a slice of building definitions, and is expected to return
// the correct slice of skyline-defining points.
func Solve(data []internal.Building) ([]image.Point, error) {
	edges := make([]image.Point, 0, len(data)*2)
	for _, b := range data {
		edges = append(edges, image.Pt(b.Sides[0], b.Height))
		edges = append(edges, image.Pt(b.Sides[1], 0))
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].X < edges[j].X
	})

	res := make([]image.Point, 0, 2*len(edges))
	open := make([]image.Point, 0, len(data))

	var cur image.Point
	res = append(res, cur)
	log.Printf("@%v", cur)

	for i, e := range edges {
		if e.Y < 0 {
			panic(fmt.Sprintf("invalid edge[%v] point %v", i, e))
		}
		if e.Y > 0 {
			if cur.X != e.X {
				cur.X = e.X
				res = append(res, cur)
				log.Printf("@%v", cur)
			}
			if cur.Y != e.Y {
				cur.Y = e.Y
				res = append(res, cur)
				log.Printf("@%v", cur)
			}
			open = append(open, e)

		} else if e.Y == 0 {
			for i := len(open) - 1; i >= 0; i-- {
				if open[i].Eq(e) {
					copy(open[i:], open[i+1:])
					open = open[:len(open)-1]
					break
				}
			}
			var level int
			for _, op := range open {
				if op.Y > level {
					level = op.Y
				}
			}

			cur.X = e.X
			if cur.Y != level {
				cur.Y = level
				res = append(res, cur)
				log.Printf("@%v", cur)
			}
		}

	}
	if len(open) > 0 {
		cur.Y = 0
		res = append(res, cur)
		log.Printf("@%v", cur)
	}

	return res, nil
}
