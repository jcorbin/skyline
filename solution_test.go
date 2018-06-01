package main_test

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/jcorbin/skyline"
	"github.com/jcorbin/skyline/internal"
)

func TestSolve_basics(t *testing.T) {
	for _, tc := range []struct {
		name     string
		data     []internal.Building
		expected []image.Point
	}{
		{
			name:     "empty data",
			data:     nil,
			expected: nil,
		},

		{
			name: "single in the middle",
			/* 0 2 4 6
			 *   |-|
			 *   | |
			 * __| |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 4}, Height: 3},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 3},
				{X: 4, Y: 3},
				{X: 4, Y: 0},
			},
		},

		{
			name: "twin towers",
			/* 0 2 4 6 8 a
			 *   |-| |-|
			 *   | | | |
			 * __| |_| |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 4}, Height: 3},
				{Sides: [2]int{6, 8}, Height: 3},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 3},
				{X: 4, Y: 3},
				{X: 4, Y: 0},

				{X: 6, Y: 0},
				{X: 6, Y: 3},
				{X: 8, Y: 3},
				{X: 8, Y: 0},
			},
		},

		{
			name: "joined towers",
			/* 0 2 4 6 8 a c e
			 *   |---| |---|
			 *   |   | |   |
			 *   |   | |   |
			 *   | ..|-|.. |
			 * __| .     . |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 6}, Height: 5},
				{Sides: [2]int{8, 12}, Height: 5},
				{Sides: [2]int{4, 10}, Height: 3},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 5},
				{X: 6, Y: 5},

				{X: 6, Y: 3},
				{X: 8, Y: 3},

				{X: 8, Y: 5},
				{X: 12, Y: 5},
				{X: 12, Y: 0},
			},
		},

		{
			name: "L",
			/*
			 * 0 2 4 6 8 a c e
			 *   |---|
			 *   | ..|__
			 * __| . . |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 6}, Height: 3},
				{Sides: [2]int{4, 8}, Height: 1},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 3},
				{X: 6, Y: 3},

				{X: 6, Y: 1},

				{X: 8, Y: 1},
				{X: 8, Y: 0},
			},
		},

		{
			name: "stair",
			/*
			 * 0 2 4 6 8 a c e
			 *   |---|
			 *   |   |
			 *   | ..|---|
			 *   | . . ..|__
			 * __| . . . . |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 6}, Height: 5},
				{Sides: [2]int{4, 10}, Height: 3},
				{Sides: [2]int{8, 12}, Height: 1},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 5},
				{X: 6, Y: 5},

				{X: 6, Y: 3},
				{X: 10, Y: 3},
				{X: 10, Y: 1},

				{X: 12, Y: 1},
				{X: 12, Y: 0},
			},
		},

		{
			name: "mirror stair",
			/*
			 * 0 2 4 6 8 a c e
			 *         |---|
			 *         |   |
			 *     |---|.. |
			 *   __|.. . . |
			 * __| . . . . |__
			 */
			data: []internal.Building{
				{Sides: [2]int{2, 6}, Height: 1},
				{Sides: [2]int{4, 10}, Height: 3},
				{Sides: [2]int{8, 12}, Height: 5},
			},
			expected: []image.Point{
				{X: 2, Y: 0},
				{X: 2, Y: 1},
				{X: 4, Y: 1},

				{X: 4, Y: 3},

				{X: 8, Y: 3},
				{X: 8, Y: 5},

				{X: 12, Y: 5},
				{X: 12, Y: 0},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			points, err := Solve(append([]internal.Building(nil), tc.data...))
			require.NoError(t, err, "expected Solve() to not fail")
			assert.Equal(t, tc.expected, points, "expected output points")
		})
	}
}

func TestSolve_gen(t *testing.T) {
	for _, tc := range []genTestCase{
		{
			seed: 0,
			w:    16,
			h:    10,
			n:    1,
		},
		{
			seed: 0,
			w:    16,
			h:    10,
			n:    2,
		},
		{
			seed: 0,
			w:    16,
			h:    10,
			n:    3,
		},
		{
			seed: 0,
			w:    16,
			h:    32,
			n:    8,
		},
	} {
		t.Run(tc.desc(), tc.run)
	}
}

type genTestCase struct {
	seed    int64
	w, h, n int
}

func (tc genTestCase) desc() string {
	return fmt.Sprintf("seed=%v w=%v h=%v n=%v", tc.seed, tc.w, tc.h, tc.n)
}

func (tc genTestCase) do(sol func([]internal.Building) []image.Point) genTestCaseRun {
	tr := genTestCaseRun{genTestCase: tc}
	tr.rng = rand.New(rand.NewSource(tr.seed))
	tr.data = internal.GenBuildings(tr.rng, tr.w, tr.h, tr.n)
	tr.points, tr.err = sol(append([]internal.Building(nil), tr.data...))
	return tr
}

type genTestCaseRun struct {
	genTestCase
	rng    *rand.Rand
	data   []internal.Building
	points []image.Point
	err    error
}

func (tr genTestCaseRun) expectedPlot() *image.Gray {
	expected := image.NewGray(image.Rect(0, 0, tr.w+2, tr.h+2))
	plotBuildings(expected, tr.data, 0x80)
	return expected
}

func (tr genTestCaseRun) actualPlot() (*image.Gray, error) {
	actual := image.NewGray(image.Rect(0, 0, tr.w+2, tr.h+2))
	if err := plotSkyline(actual, tr.points, 0x80); err != nil {
		return nil, err
	}
	return actual, nil
}

func (tc genTestCase) run(t *testing.T) {
	tr := tc.do(Solve)
	require.NoError(t, tr.err, "expected Solve() to not fail")
	actual, err := tr.actualPlot()
	require.NoError(t, err, "unable to plot skyline")
	expected := tr.expectedPlot()
	floodFill(actual, actual.Rect.Max.Sub(image.Pt(1, 1)), 0x00, 0xff)
	floodFill(expected, expected.Rect.Max.Sub(image.Pt(1, 1)), 0x00, 0xff)
	erase(actual, 0x80)
	erase(expected, 0x80)
	if !assert.Equal(t, expected, actual) {
		dumpRunes := map[uint8]rune{0x00: ' ', 0x80: '.', 0xff: '#'}

		t.Logf("building data: %v", tr.data)
		t.Logf("solution points: %v", tr.points)

		t.Logf("building plot:\n%v", strings.Join(dump(tr.expectedPlot(), dumpRunes), "\n"))
		t.Logf("expected sky:\n%v", strings.Join(dump(expected, dumpRunes), "\n"))

		skylinePlot, err := tr.actualPlot()
		require.NoError(t, err, "unable to plot skyline")
		t.Logf("skyline solution:\n%v", strings.Join(dump(skylinePlot, dumpRunes), "\n"))
		t.Logf("actual sky:\n%v", strings.Join(dump(actual, dumpRunes), "\n"))
	}
}

func dump(gr *image.Gray, tr map[uint8]rune) []string {
	res := make([]string, gr.Rect.Dy(), gr.Rect.Dy()+1)
	row := make([]rune, gr.Rect.Dx())
	for y := 0; y < len(res); y++ {
		for x := 0; x < len(row); x++ {
			if r, def := tr[gr.GrayAt(x, y).Y]; def {
				row[x] = r
			} else {
				row[x] = '?'
			}
		}
		res[len(res)-1-y] = fmt.Sprintf("% 3d: %s", y, string(row))
	}
	for x := 0; x < len(row); x++ {
		row[x] = rune(strconv.FormatInt(int64(x)%10, 10)[0])
	}
	res = append(res, fmt.Sprintf("     %s", string(row)))
	return res
}

func plotBuildings(gr *image.Gray, bs []internal.Building, val uint8) {
	for _, b := range bs {
		plotHLine(gr, b.Sides[0], b.Sides[1], 0, val)
		plotHLine(gr, b.Sides[0], b.Sides[1], b.Height, val)
		plotVLine(gr, b.Sides[0], 0, b.Height, val)
		plotVLine(gr, b.Sides[1], 0, b.Height, val)
	}
}

func plotSkyline(gr *image.Gray, points []image.Point, val uint8) error {
	if len(points) == 0 {
		return nil
	}
	errSkylinePoint := errors.New("skyline point must share exactly one component with prior")
	cur := points[0]
	for _, pt := range points[1:] {
		if pt.Eq(cur) {
			return errSkylinePoint
		}
		if pt.X == cur.X {
			plotVLine(gr, cur.X, cur.Y, pt.Y, val)
			cur.Y = pt.Y
		} else if pt.Y == cur.Y {
			plotHLine(gr, cur.X, pt.X, cur.Y, val)
			cur.X = pt.X
		} else {
			return errSkylinePoint
		}
	}
	return nil
}

func plotHLine(gr *image.Gray, x0, x1, y int, val uint8) {
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	for x := x0; x < x1; x++ {
		gr.SetGray(x, y, color.Gray{val})
	}
}

func plotVLine(gr *image.Gray, x, y0, y1 int, val uint8) {
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	for y := y0; y < y1; y++ {
		gr.SetGray(x, y, color.Gray{val})
	}
}

func erase(gr *image.Gray, val uint8) {
	for i := 0; i < len(gr.Pix); i++ {
		if gr.Pix[i] == val {
			gr.Pix[i] = 0x00
		}
	}
}

func floodFill(gr *image.Gray, pt image.Point, where, with uint8) {
	// TODO below is most naive / slow implementation possible, probably would
	// be worth it to move to at least a scanline approach
	// q := make([]image.Point, 0, gr.Rect.Dy()+1)

	q := make([]image.Point, 0, gr.Rect.Dx()*gr.Rect.Dy())
	sanity := 2 * 5 * cap(q)
	q = append(q, pt)
	for len(q) > 0 {
		if sanity--; sanity <= 0 {
			panic(fmt.Sprintf(
				"suspect infinite loop in floodFill(%v, %v, %v, %v)",
				gr.Rect, pt, where, with,
			))
		}
		var qpt image.Point
		qpt, q = q[0], q[:copy(q, q[1:])]
		if fill(gr, qpt, where, with) {
			for _, d := range []image.Point{
				image.Pt(0, 1), image.Pt(0, -1),
				image.Pt(1, 0), image.Pt(-1, 0),
			} {
				npt := qpt.Add(d)
				if npt.In(gr.Rect) {
					q = append(q, npt)
				}
			}
		}
	}
}

func fill(gr *image.Gray, pt image.Point, where, with uint8) bool {
	if gr.GrayAt(pt.X, pt.Y).Y == where {
		gr.SetGray(pt.X, pt.Y, color.Gray{with})
		return true
	}
	return false
}
