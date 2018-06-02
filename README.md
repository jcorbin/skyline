# Contemplating the "skyline" interview problem

> Given a 2d description of a set of buildings, compute their skyline profile.

See the [master](../../tree/master) branch for more on the basic problem
description.

# Go Setup

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

# TODO

- verify skyline plot fail code path
- maybe shortcircuit not implemented error, dont run cases
- unify plotting layer around gray
- add benchmark infra
- turn the perf crank
