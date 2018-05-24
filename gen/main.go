package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
)

func main() {
	var (
		widthFlag  = flag.Int("w", 16, "width of field")
		heightFlag = flag.Int("h", 32, "height of field")
		countFlag  = flag.Int("n", 8, "number of buildings")
		seedFlag   = flag.Int64("s", 0, "rand seed (default: random)")
	)
	flag.Parse()
	if err := run(
		os.Stdout,
		*seedFlag,
		*widthFlag,
		*heightFlag,
		*countFlag,
	); err != nil {
		log.Fatalln(err)
	}
}

func run(
	out io.Writer,
	seed int64,
	w, h, n int,
) error {
	rng := rand.New(rand.NewSource(seed))

	if _, err := fmt.Fprintf(out, "x1,x2,h\n"); err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		x1 := rng.Intn(w)
		x2 := rng.Intn(w)
		if x2 < x1 {
			x1, x2 = x2, x1
		}
		h := rng.Intn(h-1) + 1
		if _, err := fmt.Fprintf(out, "%d,%d,%d\n", x1, x2, h); err != nil {
			return err
		}
	}
	return nil
}
