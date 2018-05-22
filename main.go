package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jcorbin/skyline/internal"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		log.Fatalln(err)
	}
}

func run(in io.Reader, out io.Writer) error {
	data, err := internal.ScanBuildings(in)
	if err != nil {
		return err
	}
	points, err := Solve(data)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "x,y\n"); err != nil {
		return err
	}
	for _, pt := range points {
		if _, err := fmt.Fprintf(out, "%d,%d\n", pt.X, pt.Y); err != nil {
			return err
		}
	}
	return nil
}
