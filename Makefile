all: bin/gen bin/display

NMIN=0
NMAX=1024
NSTEPS=32
TESTFLAGS=-gen.nmin $(NMIN) -gen.nmax $(NMAX) -gen.nsteps $(NSTEPS)

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
