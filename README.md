# Contemplating the "skyline" interview problem

> Given a 2d description of a set of buildings, compute their skyline profile.

Building description is fairly straight forward, there are some basic choices:
- rectangle defined by two points
- left/right X values with a height

For "the skyline" however, desired outcome becomes less clear:
- Do you just want to trace its outline? If so a list of `<X, Y>` points that
  describe each corner works.
- But if you wanted to fill it in, you might instead want directly compute a
  rectangle strip; such strip itself has many representation questions.
- Furthermore, if you were really drawing this stuff (say targeting the GL
  api), you'd could directly triangulate the skyline (compute a triangle
  strip).

# Go Setup

This branch contains support and a starting point for go solutions (return to
the [master](../../tree/master) branch).

Some convenience utilities:
- a basic problem generator in the [gen](./gen) package
- a problem displayer in the [display](./display) package
- such utilities can be built with a `make all`, then run from `./bin/`

There's a [main.go](main harness) provided, as well as some basic [test
cases](solution_test.go).

Solution implementation can start by implementing in `solve.go` the stubbed in:
```golang
func Solve(data []internal.Building) ([]image.Point, error)
```

Solution may be ran over random data, for example:
```shell
$ ./bin/gen | go run main.go solution.go
```
