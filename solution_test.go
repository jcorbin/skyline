package main_test

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/jcorbin/skyline"
	"github.com/jcorbin/skyline/internal"
)

var staticTestCases = []testCase{
	{
		name:   "empty data",
		data:   nil,
		points: nil,
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
		points: []image.Point{
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
		points: []image.Point{
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
		points: []image.Point{
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
		points: []image.Point{
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
		points: []image.Point{
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
		points: []image.Point{
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
}

var genTestCases = []testCase{
	{
		seed: 0,
		w:    16,
		h:    10,
	},
	{
		seed: 0,
		w:    32,
		h:    32,
	},
	{
		seed: 0,
		w:    64,
		h:    64,
	},
}

func TestSolve(t *testing.T) {
	if _, err := Solve(nil); err != nil {
		t.Logf("Solve() failed unequivocally: %v", err)
		t.Fail()
		return
	}
	for _, tc := range staticTestCases {
		t.Run(tc.String(), tc.run(Solve).runTest)
	}
	for _, tc := range genTestCases {
		t.Run(tc.String(), tc.run(Solve).runTest)
	}
}

func TestSolver_Solve(t *testing.T) {
	var sol Solver
	if _, err := sol.Solve(nil); err != nil {
		t.Logf("sol.Solve() failed unequivocally: %v", err)
		t.Fail()
		return
	}
	for _, tc := range staticTestCases {
		t.Run(tc.String(), tc.run(sol.Solve).runTest)
	}
	for _, tc := range genTestCases {
		t.Run(tc.String(), tc.run(sol.Solve).runTest)
	}
}

func BenchmarkSolve(b *testing.B) {
	for _, tc := range genTestCases {
		b.Run(tc.String(), tc.run(Solve).runBench)
	}
}

type testCase struct {
	name   string
	data   []internal.Building
	points []image.Point

	// generative parameters
	seed    int64
	w, h, n int
}

type testCaseRun struct {
	sol func([]internal.Building) ([]image.Point, error)
	testCase

	gen bool
	rng *rand.Rand

	points []image.Point

	buildingPlot, skylinePlot *image.Gray
	expectedSky, actualSky    *image.Gray
}

func (tc testCase) isGen() bool {
	if tc.data != nil {
		return false
	}
	return tc.w > 0 && tc.h > 0
}

func (tc testCase) String() string {
	if tc.name != "" {
		// manually named case
		return tc.name
	}
	if tc.isGen() {
		if tc.n == 0 {
			return fmt.Sprintf("genTest<seed=%v w=%v h=%v>", tc.seed, tc.w, tc.h)
		}
		return fmt.Sprintf("genTest<seed=%v w=%v h=%v n=%v>", tc.seed, tc.w, tc.h, tc.n)
	}
	return fmt.Sprintf("staticTest<data=%v points=%v>", tc.data, tc.points)
}

func (tc testCase) run(sol func([]internal.Building) ([]image.Point, error)) testCaseRun {
	return testCaseRun{
		testCase: tc,
		sol:      sol,
		gen:      tc.isGen(),
	}
}

func (tr *testCaseRun) solve(data []internal.Building) (err error) {
	tr.points, err = tr.sol(data)
	return err
}

func (tr testCaseRun) runTest(t *testing.T) {
	defer setupTestLogOutput(t).restore(os.Stderr)
	if !tr.gen {
		tr.doStaticTest(t)
	} else if tr.n != 0 {
		tr.doGenTest(t)
	} else {
		tr.doGenSearchTest(t)
	}
}

func (tr testCaseRun) doStaticTest(t *testing.T) {
	data := append([]internal.Building(nil), tr.data...)
	require.NoError(t, tr.solve(data), "expected solution to not fail")
	if !assert.Equal(t, tr.testCase.points, tr.points, "expected output points") {
		if err := tr.buildPlots(); err != nil {
			t.Logf("unable to plot skyline: %v", err)
		} else {
			tr.logDebugInfo(t.Logf)
		}
	}
}

func (tr testCaseRun) doGenTest(t *testing.T) {
	tr.rng = rand.New(rand.NewSource(tr.seed))
	tr.data = internal.GenBuildings(tr.rng, tr.w, tr.h, tr.n)
	data := append([]internal.Building(nil), tr.data...)
	require.NoError(t, tr.solve(data), "expected solution to not fail")
	require.NoError(t, tr.buildPlots(), "unable to plot skyline")
	if !grayEQ(tr.expectedSky, tr.actualSky) {
		t.Fail()
		tr.logDebugInfo(t.Logf)
	}
}

func (tr testCaseRun) doGenSearchTest(t *testing.T) (pass bool) {
	const (
		min  = 1
		max  = 1024
		step = 128
	)

	if tr.n = min; !t.Run(tr.String(), tr.doGenTest) {
		return false
	}

	pass = true
	for tr.n = max; tr.n > min; tr.n -= step {
		if !t.Run(tr.String(), tr.doGenTest) {
			pass = false
			break
		}
	}
	if pass {
		return true
	}

	sanity := max - min
	for n := min; tr.n-n > 1; {
		sanity--
		require.True(t, sanity > 0, "search looping infinitely")
		lastN := tr.n
		tr.n = lastN/2 + n/2
		if t.Run(tr.String(), tr.doGenTest) {
			n, tr.n = tr.n, lastN
		}
	}
	t.Logf("found minimal failure case in %v", tr)
	return false
}

func (tr testCaseRun) runBench(b *testing.B) {
	defer setupTestLogOutput(b).restore(os.Stderr)
	if !tr.gen {
		tr.doStaticBench(b)
	} else if tr.n != 0 {
		tr.doGenBench(b)
	} else {
		tr.doGenScaleBench(b)
	}
}

func (tr testCaseRun) doStaticBench(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := append([]internal.Building(nil), tr.data...)
		require.NoError(b, tr.solve(data), "expected solution to not fail")
	}
}

func (tr testCaseRun) doGenBench(b *testing.B) {
	tr.rng = rand.New(rand.NewSource(tr.seed))
	tr.data = internal.GenBuildings(tr.rng, tr.w, tr.h, tr.n)
	for i := 0; i < b.N; i++ {
		data := append([]internal.Building(nil), tr.data...)
		require.NoError(b, tr.solve(data), "expected solution to not fail")
	}
}

func (tr testCaseRun) doGenScaleBench(b *testing.B) {
	const (
		min  = 0
		max  = 1024
		step = 32
	)
	for tr.n = min; tr.n < max; tr.n += step {
		b.Run(tr.String(), tr.doGenBench)
	}
}

func (tr testCaseRun) logDebugInfo(logf func(string, ...interface{})) {
	logf("building data: %v", tr.data)
	logf("solution points: %v", tr.points)
	dumpRunes := map[uint8]rune{0x00: ' ', 0x80: '.', 0xff: '#'}
	logf("plots:\n%s", strings.Join(sideBySide(
		"building boxes", "skyline",
		dump(tr.buildingPlot, dumpRunes),
		dump(tr.skylinePlot, dumpRunes),
	), "\n"))
	logf("skies:\n%s", strings.Join(sideBySide(
		"expected", "actual",
		dump(tr.expectedSky, dumpRunes),
		dump(tr.actualSky, dumpRunes),
	), "\n"))
}

func (tr *testCaseRun) buildPlots() error {
	oob := image.Pt(tr.w+1, tr.h+1) // out-of-bounds fill starting point
	tr.buildingPlot = image.NewGray(image.Rect(0, 0, oob.X+1, oob.Y+1))
	tr.skylinePlot = image.NewGray(image.Rect(0, 0, oob.X+1, oob.Y+1))
	plotBuildings(tr.buildingPlot, tr.data, 0x80)
	if err := plotSkyline(tr.skylinePlot, tr.points, 0x80); err != nil {
		return err
	}
	tr.expectedSky = plot2sky(tr.buildingPlot, oob, 0x00, 0xff, 0x80)
	tr.actualSky = plot2sky(tr.skylinePlot, oob, 0x00, 0xff, 0x80)
	return nil
}

type testLogOutput struct {
	testing.TB
	priorFlags int
}

func setupTestLogOutput(tb testing.TB) testLogOutput {
	var tlo testLogOutput
	tlo.TB = tb
	tlo.priorFlags = log.Flags()
	log.SetFlags(0)
	log.SetOutput(tlo)
	return tlo
}

func (tlo testLogOutput) restore(priorOutput io.Writer) {
	log.SetOutput(priorOutput)
	log.SetFlags(tlo.priorFlags)
}

func (tlo testLogOutput) Write(p []byte) (int, error) {
	tlo.Logf("%s", p)
	return 0, nil
}

func sideBySide(aTitle, bTitle string, a, b []string) []string {
	var res []string
	if len(b) > len(a) {
		res = make([]string, 0, len(b)+1)
	} else {
		res = make([]string, 0, len(a)+1)
	}
	var aw, bw int
	for _, as := range a {
		if len(as) > aw {
			aw = len(as)
		}
	}
	for _, bs := range b {
		if len(bs) > bw {
			bw = len(bs)
		}
	}
	if aTitle != "" || bTitle != "" {
		res = append(res, fmt.Sprintf("| % *s | % *s |", aw, aTitle, bw, bTitle))
	}
	for i := 0; i < len(a) || i < len(b); i++ {
		var as, bs string
		if i < len(a) {
			as = a[i]
		}
		if i < len(b) {
			bs = b[i]
		}
		res = append(res, fmt.Sprintf("| % *s | % *s |", aw, as, bw, bs))
	}
	return res
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
	for x := x0; x <= x1; x++ {
		gr.SetGray(x, y, color.Gray{val})
	}
}

func plotVLine(gr *image.Gray, x, y0, y1 int, val uint8) {
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	for y := y0; y <= y1; y++ {
		gr.SetGray(x, y, color.Gray{val})
	}
}

func plot2sky(
	gr *image.Gray,
	fillAt image.Point, fillWhere, fillWith uint8,
	eraseWhere uint8,
) *image.Gray {
	ngr := image.NewGray(gr.Rect)
	copy(ngr.Pix, gr.Pix)

	// fill sky
	floodFill(ngr, fillAt, fillWhere, fillWith)

	// erase plot
	for i := 0; i < len(ngr.Pix); i++ {
		if ngr.Pix[i] == eraseWhere {
			ngr.Pix[i] = 0x00
		}
	}

	return ngr
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

func grayEQ(a, b *image.Gray) bool {
	if !a.Rect.Eq(b.Rect) {
		return false
	}
	for i := range a.Pix {
		if a.Pix[i] != b.Pix[i] {
			return false
		}
	}
	return true
}
