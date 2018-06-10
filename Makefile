all: bin/gen bin/display

# case generation parameters for tests and benchmarks:
# - N buildings are randomly generated using a seeded RNG
# - within a WxH world (the xH is optional, eliding it for square)
# - test at MIN and MAX; if that passes, test at each STEP between MIN and MAX;
#   once a failure is found, minimize its N value using binary search.
# - benchmark from MIN to MAX in STEPs
GENSEEDS=0
GENSIZES=16,32,64
NMIN=0
NMAX=1024
NSTEPS=32

TESTFLAGS=-gen.seeds $(GENSEEDS) -gen.sizes $(GENSIZES) -gen.nmin $(NMIN) -gen.nmax $(NMAX) -gen.nsteps $(NSTEPS)

test:
	go test -v . $(TESTFLAGS)

BENCHTIME:=100ms

ifeq ($(strip $(OUT)),)
bench:
	go test . $(TESTFLAGS) -bench . -benchmem -benchtime $(BENCHTIME) \
		-o skyline.test \
		-cpuprofile cpu.pprof \
		-memprofile mem.pprof \
		| tee bench.out
else
bench:
	mkdir -p $(OUT)
	go test . $(TESTFLAGS) -bench . -benchmem -benchtime $(BENCHTIME) \
		-o $(OUT)/skyline.test \
		-cpuprofile $(OUT)/cpu.pprof \
		-memprofile $(OUT)/mem.pprof \
		| tee $(OUT)/bench.out
endif

bin/gen: gen/*.go
	go build -o $@ $^

bin/display: display/*.go
	go build -o $@ $^
