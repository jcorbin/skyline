package internal

import (
	"math/rand"
)

// Gen erate n random buildings chosen (rather poorly!) using the given
// random generator and width/height constraints; calls f with each building,
// passing through and stopping on any non-nil error.
//
// TODO this generator could be higher "quality"
func Gen(
	rng *rand.Rand,
	w, h, n int,
	f func(b Building) error,
) error {
	for i := 0; i < n; i++ {
		x1 := rng.Intn(w)
		x2 := rng.Intn(w)
		if x2 < x1 {
			x1, x2 = x2, x1
		}
		h := rng.Intn(h-1) + 1
		if err := f(Building{[2]int{x1, x2}, h}); err != nil {
			return err
		}
	}
	return nil
}

// GenBuildings is a convenience wrapper around Gen, returning n-randomly
// generated buildings.
func GenBuildings(rng *rand.Rand, w, h, n int) (bs []Building) {
	bs = make([]Building, n)
	Gen(rng, w, h, n, func(b Building) error {
		bs = append(bs, b)
		return nil
	})
	return bs
}
