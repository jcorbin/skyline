# Contemplating the "skyline" interview problem

> Given a 2d description of a set of buildings, compute their skyline profile.

See the [master](../../tree/master) branch for more on the basic problem
description.

# Go Setup

Some convenience utilities:
- a basic problem generator in the [gen](./gen) package
- a problem displayer in the [display](./display) package
- such utilities can be built with a `make all`, then run from `./bin/`

Solution implementation can start by implementing in `solve.go` the stubbed in:
```golang
func Solve(data []internal.Building) ([]image.Point, error)
```

## Getting Started

To start a go solution branch, run `./start.bash` to create a
`go_solution_YYYYMMDD` branch; if ran from a pre-existing
`go_solution_YYYYMMDD` branch, then any prior (maybe retcon-ed) progress, is
used as a starting point.

Solution may be ran over random data, for example:
```shell
$ ./bin/gen | go run main.go solution.go
```

## Testing

There's an extensive [test suite](solution_test.go) provided:
- It has basic test cases like empty data, a single tower, two
  disjoint towers, and so on.
- If the basic test cases all pass, then a further set of more
  exhaustive stress tests are ran:
  - They use the same random data generator as the [gen](./gen) program...
  - ...to progressively test the solution with more and more buildings
  - ...and using binary search to narrow down any found failure to a minimum
    number of buildings.
  - This testing process is carried out on varying world sizes.
  - The solution fails this exhaustive / generative test if:
    - Its resulting skyline has anything other than horizontal or vertical
      lines.
    - The resulting sky isn't the same as the sky resulting from the building
      box input shape.

The tests can be run using `make test` for convenience, or directly using `go
test` if you prefer.

## Benchmarking

Once your solution passes, you can furthermore take a look at "how efficient /
fast is it?"

The [benchmark suite](solution_test.go) uses all the same cases as the test
suite above, trying to run the solution as many times as possible within an
allotted time. The generative test cases benchmark at progressively larger
numbers of buildings, from which you can extrapolate how well your solution
scales.

The benchmarks can be run using `make bench` for convenience (which is
just a wrapper around `go test -bench ...`):
- CPU profiling is enabled, and writes to `cpu.pprof`
- Memory profiling is enabled, and writes to `mem.pprof`
- The test binary is written to `skyline.test`
- The resulting benchmark stats (ns, B, and allocs /op) is written to
  `bench.out` (as well as the console).
- You may further specify `make OUT=name bench` to stoe all of the
  above files in a `name/`ed sub-directory; this supports easy
  comparison of before/after and so on.
