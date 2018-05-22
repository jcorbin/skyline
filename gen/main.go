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
	)
	flag.Parse()
	if err := run(
		os.Stdout,
		*widthFlag,
		*heightFlag,
		*countFlag,
	); err != nil {
		log.Fatalln(err)
	}
}

func run(out io.Writer, w, h, n int) error {
	if _, err := fmt.Fprintf(out, "x1,x2,h\n"); err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		x1 := rand.Intn(w)
		x2 := rand.Intn(w)
		if x2 < x1 {
			x1, x2 = x2, x1
		}
		h := rand.Intn(h)
		if _, err := fmt.Fprintf(out, "%d,%d,%d\n", x1, x2, h); err != nil {
			return err
		}
	}
	return nil
}
