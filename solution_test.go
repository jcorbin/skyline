package main_test

import (
	"bufio"
	"bytes"
	"flag"
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

var (
	genMin   = 0
	genMax   = 1024
	genStep  = 32
	genSeeds []int64
	genSizes = []image.Point{
		image.Pt(16, 16),
		image.Pt(32, 32),
		image.Pt(64, 64),
	}
)

type _genSeeds struct{}

func (gs _genSeeds) String() string {
	parts := make([]string, len(genSeeds))
	for i, seed := range genSeeds {
		parts[i] = strconv.FormatInt(seed, 10)
	}
	return strings.Join(parts, ",")
}

func (gs _genSeeds) Set(s string) error {
	genSeeds = genSeeds[:0]
	for _, part := range strings.Split(s, ",") {
		seed, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return err
		}
		genSeeds = append(genSeeds, seed)
	}
	return nil
}

type _genSteps struct{}

func (gs _genSteps) String() string {
	w := genMax - genMin
	n := (w + genStep - 1) / genStep
	return strconv.Itoa(n)
}

func (gs _genSteps) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	w := genMax - genMin
	genStep = w / n
	return nil
}

type sizesFlag struct {
	sizes *[]image.Point
}

func (ss sizesFlag) String() string {
	var sizes []image.Point
	if ss.sizes != nil {
		sizes = *ss.sizes
	}
	parts := make([]string, len(sizes))
	for i, sz := range sizes {
		if sz.X == sz.Y {
			parts[i] = fmt.Sprintf("%v", sz.X)
		} else {
			parts[i] = fmt.Sprintf("%vx%v", sz.X, sz.Y)
		}
	}
	return strings.Join(parts, ",")
}

func (ss sizesFlag) Set(s string) (err error) {
	sizes := (*ss.sizes)[:0]
	for _, part := range strings.Split(s, ",") {
		var sz image.Point
		if i := strings.Index(part, "x"); i > 0 {
			sz.X, err = strconv.Atoi(part[:i])
			if err == nil {
				sz.Y, err = strconv.Atoi(part[i+1:])
			}
		} else {
			sz.X, err = strconv.Atoi(part)
			sz.Y = sz.X
		}
		if err != nil {
			return err
		}
		sizes = append(sizes, sz)
	}
	*ss.sizes = sizes
	return err
}

func init() {
	flag.Var(_genSeeds{}, "gen.seeds", "custom seed(s) for generating test data")
	flag.Var(sizesFlag{sizes: &genSizes}, "gen.sizes", "world sizes for generating test data")
	flag.IntVar(&genMin, "gen.nmin", genMin,
		"minimum N value for generative tests and benchmarks")
	flag.IntVar(&genMax, "gen.nmax", genMax,
		"maximum N value for generative tests and benchmarks")
	flag.IntVar(&genStep, "gen.nstep", genStep,
		"linear N step size for generative tests and benchmarks")
	flag.Var(_genSteps{}, "gen.nsteps", "convenience for setting -gen.nstep")
}

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

func startGen(next func() (testCase, bool)) (func() (testCase, bool), testCase, bool) {
	tc, ok := next()
	return next, tc, ok
}

func genCases() func() (testCase, bool) {
	if len(genSeeds) == 0 {
		return genSizeCases(0, genSizes...)
	}
	i := 0
	var nextSize func() (testCase, bool)
	return func() (tc testCase, ok bool) {
		if nextSize != nil {
			tc, ok = nextSize()
			if !ok && i < len(genSeeds) {
				nextSize = nil
			}
		}
		if nextSize == nil && i < len(genSeeds) {
			seed := genSeeds[i]
			i++
			nextSize = genSizeCases(seed, genSizes...)
			tc, ok = nextSize()
		}
		return tc, ok
	}
}

func genSizeCases(seed int64, sizes ...image.Point) func() (testCase, bool) {
	i := 0
	return func() (testCase, bool) {
		if i < len(sizes) {
			sz := sizes[i]
			i++
			return testCase{
				seed: seed,
				w:    sz.X,
				h:    sz.Y,
			}, true
		}
		return testCase{}, false
	}
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
	if !t.Failed() {
		for next, tc, ok := startGen(genCases()); ok; tc, ok = next() {
			t.Run(tc.String(), tc.run(Solve).runTest)
		}
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
	if !t.Failed() {
		for next, tc, ok := startGen(genCases()); ok; tc, ok = next() {
			t.Run(tc.String(), tc.run(sol.Solve).runTest)
		}
	}
}

func BenchmarkSolver_Solve(b *testing.B) {
	var sol Solver
	for _, tc := range staticTestCases {
		b.Run(tc.String(), tc.run(sol.Solve).runBench)
	}
	for next, tc, ok := startGen(genCases()); ok; tc, ok = next() {
		b.Run(tc.String(), tc.run(sol.Solve).runBench)
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
	bufa, bufb, bufc          *bytes.Buffer
}

func (tc testCase) isGen() bool {
	if tc.data != nil {
		return false
	}
	return tc.w > 0 && tc.h > 0
}

func (tc testCase) maxPoint() image.Point {
	pt := image.Pt(tc.w, tc.h)
	if pt == image.ZP {
		for _, b := range tc.data {
			if pt.X < b.Sides[1] {
				pt.X = b.Sides[1]
			}
			if pt.Y < b.Height {
				pt.Y = b.Height
			}
		}
	}
	return pt
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
		bufa:     bytes.NewBuffer(nil),
		bufb:     bytes.NewBuffer(nil),
		bufc:     bytes.NewBuffer(nil),
	}
}

func (tr *testCaseRun) solve(data []internal.Building) (err error) {
	tr.points, err = tr.sol(data)
	return err
}

func (tr testCaseRun) runTest(t *testing.T) {
	defer setupTestLogOutput(t).restore(os.Stderr)
	if !tr.gen {
		tr.doTest(t)
	} else if tr.n != 0 {
		tr.doGenTest(t)
	} else {
		tr.doGenSearchTest(t)
	}
}

func (tr testCaseRun) doGenTest(t *testing.T) {
	tr.rng = rand.New(rand.NewSource(tr.seed))
	tr.data = internal.GenBuildings(tr.rng, tr.w, tr.h, tr.n)
	tr.doTest(t)
}

func (tr testCaseRun) doTest(t *testing.T) {
	data := append([]internal.Building(nil), tr.data...)
	require.NoError(t, tr.solve(data), "expected solution to not fail")
	if !assert.NoError(t, plotSkyline(nil, tr.points, 0x00), "expected a valid skyline") {
		t.Logf("building data: %v", tr.data)
		t.Logf("solution points: %v", tr.points)
		return
	}

	if !tr.gen {
		assert.Equal(t, tr.testCase.points, tr.points, "expected output points")
	} else if assert.NoError(t, tr.buildPlots(), "unable to plot skyline") {
		assert.True(t, grayEQ(tr.expectedSky, tr.actualSky), "expected same resulting sky")
	}

	if t.Failed() {
		t.Logf("building data: %v", tr.data)
		t.Logf("solution points: %v", tr.points)
		assert.NoError(t, tr.buildPlots(), "unable to plot skyline")
		tr.dumpPlots(t.Logf)
	}
}

func (tr testCaseRun) doGenSearchTest(t *testing.T) (pass bool) {
	if tr.n = genMin; !t.Run(tr.String(), tr.doGenTest) {
		return false
	}

	pass = true
	for tr.n = genMax; tr.n > genMin; tr.n -= genStep {
		if !t.Run(tr.String(), tr.doGenTest) {
			pass = false
			break
		}
	}
	if pass {
		return true
	}

	sanity := genMax - genMin
	for n := genMin; tr.n-n > 1; {
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
		tr.doBench(b)
	} else if tr.n != 0 {
		tr.doGenBench(b)
	} else {
		tr.doGenScaleBench(b)
	}
}

func (tr testCaseRun) doGenBench(b *testing.B) {
	tr.rng = rand.New(rand.NewSource(tr.seed))
	tr.data = internal.GenBuildings(tr.rng, tr.w, tr.h, tr.n)
	tr.doBench(b)
}

func (tr testCaseRun) doBench(b *testing.B) {
	data := make([]internal.Building, len(tr.data))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(data, tr.data)
		if err := tr.solve(data); err != nil {
			b.Logf("solution failed with unexpected error: %v", err)
			b.FailNow()
		}
	}
}

func (tr testCaseRun) doGenScaleBench(b *testing.B) {
	for tr.n = genMin; tr.n <= genMax; tr.n += genStep {
		b.Run(tr.String(), tr.doGenBench)
	}
}

func (tr testCaseRun) dumpPlots(logf func(string, ...interface{})) {
	if tr.buildingPlot != nil && tr.skylinePlot != nil {
		tr.bufc.WriteString("plots:\n")
		dump(tr.bufa, tr.buildingPlot, dumpRunes)
		dump(tr.bufb, tr.skylinePlot, dumpRunes)
		sideBySide(tr.bufc, "building boxes", "skyline", tr.bufa, tr.bufb)
		logf(tr.bufc.String())
		tr.bufa.Reset()
		tr.bufb.Reset()
		tr.bufc.Reset()
	}
	if tr.expectedSky != nil && tr.actualSky != nil {
		tr.bufc.WriteString("skies:\n")
		dump(tr.bufa, tr.expectedSky, dumpRunes)
		dump(tr.bufb, tr.actualSky, dumpRunes)
		sideBySide(tr.bufc, "expected", "actual", tr.bufa, tr.bufb)
		logf(tr.bufc.String())
		tr.bufa.Reset()
		tr.bufb.Reset()
		tr.bufc.Reset()
	}
}

func (tr *testCaseRun) buildPlots() error {
	if tr.buildingPlot != nil || tr.skylinePlot != nil {
		return nil
	}
	oob := tr.maxPoint().Add(image.Pt(1, 1)) // out-of-bounds fill starting point
	box := image.Rectangle{Max: oob.Add(image.Pt(1, 1))}
	tr.buildingPlot = image.NewGray(box)
	tr.skylinePlot = image.NewGray(box)
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

func scanMaxLineLen(r io.Reader) (lines, maxLen int) {
	for sc := bufio.NewScanner(r); sc.Scan(); {
		lines++
		if n := len(sc.Bytes()); maxLen < n {
			maxLen = n
		}
	}
	return lines, maxLen
}

func sideBySide(buf *bytes.Buffer, aTitle, bTitle string, a, b *bytes.Buffer) {
	al, aw := scanMaxLineLen(bytes.NewReader(a.Bytes()))
	bl, bw := scanMaxLineLen(bytes.NewReader(b.Bytes()))
	if aw < len(aTitle) {
		aw = len(aTitle)
	}
	if bw < len(bTitle) {
		bw = len(bTitle)
	}
	nl := al
	if nl < bl {
		nl = bl
	}

	w := 2 + aw + 3 + bw + 2 // "| " + A[i] + " | " + B[i] + " |"
	buf.Grow(nl * w)

	if aTitle != "" || bTitle != "" {
		buf.WriteString("| ")
		for i := len(aTitle); i < aw; i++ {
			buf.WriteByte(' ')
		}
		buf.WriteString(aTitle)
		buf.WriteString(" | ")
		for i := len(bTitle); i < bw; i++ {
			buf.WriteByte(' ')
		}
		buf.WriteString(bTitle)
		buf.WriteString(" |\n")
	}

	as := bufio.NewScanner(bytes.NewReader(a.Bytes()))
	bs := bufio.NewScanner(bytes.NewReader(b.Bytes()))

	for {
		var atok, btok []byte
		if as.Scan() {
			atok = as.Bytes()
		}
		if bs.Scan() {
			btok = bs.Bytes()
		}
		if len(atok) == 0 && len(btok) == 0 {
			break
		}
		buf.WriteString("| ")
		for i := len(atok); i < aw; i++ {
			buf.WriteByte(' ')
		}
		buf.Write(atok)
		buf.WriteString(" | ")
		for i := len(btok); i < bw; i++ {
			buf.WriteByte(' ')
		}
		buf.Write(btok)
		buf.WriteString(" |\n")
	}
}

func dump(buf *bytes.Buffer, gr *image.Gray, tr map[uint8]rune) {
	dx, dy := gr.Rect.Dx(), gr.Rect.Dy()
	buf.Grow((5 + dx + 1) * (dy + 1))
	for y := dy - 1; y >= 0; y-- {
		fmt.Fprintf(buf, "% 3d: ", y)
		for x := 0; x < dx; x++ {
			if r, def := tr[gr.GrayAt(x, y).Y]; def {
				buf.WriteRune(r)
			} else {
				buf.WriteByte('?')
			}
		}
		buf.WriteByte('\n')
	}
	buf.WriteString("     ")
	for x := 0; x < dx; x++ {
		buf.WriteByte(strconv.FormatInt(int64(x)%10, 10)[0])
	}
	buf.WriteByte('\n')
}

func plotBuildings(gr *image.Gray, bs []internal.Building, val uint8) {
	if gr == nil {
		return
	}
	for _, b := range bs {
		plotHLine(gr, b.Sides[0], b.Sides[1], 0, val)
		plotHLine(gr, b.Sides[0], b.Sides[1], b.Height, val)
		plotVLine(gr, b.Sides[0], 0, b.Height, val)
		plotVLine(gr, b.Sides[1], 0, b.Height, val)
	}
}

func plotSkyline(gr *image.Gray, points []image.Point, val uint8) error {
	const (
		dirNone = iota
		dirVert
		dirHoriz
	)

	if len(points) == 0 {
		return nil
	}

	last, lastDir, cur := image.ZP, dirNone, points[0]
	for i := 1; i < len(points); i++ {
		pt := points[i]
		if pt.Eq(cur) {
			return fmt.Errorf("skyline contains duplicate point [%v]=%v", i, pt)
		}
		dir := dirNone
		if pt.X == cur.X {
			dir = dirVert
			plotVLine(gr, cur.X, cur.Y, pt.Y, val)
			cur.Y = pt.Y
		} else if pt.Y == cur.Y {
			dir = dirHoriz
			plotHLine(gr, cur.X, pt.X, cur.Y, val)
			cur.X = pt.X
		} else {
			return fmt.Errorf("skyline contains diagonal line from [%v]=%v to [%v]=%v", i-1, cur, i, pt)
		}
		if dir == lastDir {
			return fmt.Errorf("skyline contains co-linear points through {[%v]=%v [%v]=%v [%v]=%v}",
				i-2, points[i-2],
				i-1, last,
				i, cur,
			)
		}
		last, lastDir = cur, dir
	}
	return nil
}

func plotHLine(gr *image.Gray, x0, x1, y int, val uint8) {
	if gr == nil {
		return
	}
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	for x := x0; x <= x1; x++ {
		gr.SetGray(x, y, color.Gray{val})
	}
}

func plotVLine(gr *image.Gray, x, y0, y1 int, val uint8) {
	if gr == nil {
		return
	}
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	for y := y0; y <= y1; y++ {
		gr.SetGray(x, y, color.Gray{val})
	}
}

var dumpRunes = map[uint8]rune{0x00: ' ', 0x80: '.', 0xff: '#'}

func plot2sky(
	gr *image.Gray,
	fillAt image.Point, fillWhere, fillWith uint8,
	eraseWhere uint8,
) *image.Gray {
	if gr == nil {
		return nil
	}
	ngr := image.NewGray(gr.Rect)
	copy(ngr.Pix, gr.Pix)

	// debug dump space
	bufa := bytes.NewBuffer(nil)
	bufb := bytes.NewBuffer(nil)
	bufc := bytes.NewBuffer(nil)

	// fill sky
	fmt.Fprintf(bufc, "flood fill sky @%v where %v with %v in:\n", fillAt, fillWhere, fillWith)
	dump(bufa, ngr, dumpRunes)
	floodFill(ngr, fillAt, fillWhere, fillWith)
	dump(bufb, ngr, dumpRunes)
	sideBySide(bufc, "before", "after", bufa, bufb)
	log.Printf(bufc.String())

	// erase plot
	bufa.Reset()
	bufb.Reset()
	bufc.Reset()
	fmt.Fprintf(bufc, "erase %v:\n", eraseWhere)
	dump(bufa, ngr, dumpRunes)
	for i := 0; i < len(ngr.Pix); i++ {
		if ngr.Pix[i] == eraseWhere {
			ngr.Pix[i] = 0x00
		}
	}
	dump(bufb, ngr, dumpRunes)
	sideBySide(bufc, "before", "after", bufa, bufb)
	log.Printf(bufc.String())

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
		qpt := q[0]
		q = q[:copy(q, q[1:])]
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
