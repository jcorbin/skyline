package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"

	"github.com/jcorbin/skyline/internal"
)

func main() {
	var (
		widthFlag  = flag.Int("w", 16, "width of field")
		heightFlag = flag.Int("h", 32, "height of field")
		countFlag  = flag.Int("n", 8, "number of buildings")
		seedFlag   = flag.Int64("s", 0, "rand seed (default: random)")
	)
	flag.Parse()
	_, err := fmt.Printf("x1,x2,h\n")
	if err == nil {
		err = internal.Gen(
			rand.New(rand.NewSource(*seedFlag)),
			*widthFlag,
			*heightFlag,
			*countFlag,
			func(b internal.Building) error {
				_, err := fmt.Printf("%d,%d,%d\n", b.Sides[0], b.Sides[1], b.Height)
				return err
			})
	}
	if err != nil {
		log.Fatalln(err)
	}
}
