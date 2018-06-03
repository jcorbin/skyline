all: bin/gen bin/display

test:
	go test -v .

BENCHTIME:=100ms

bench:
	go test . -bench . -benchmem -benchtime $(BENCHTIME) | tee bench.out

ifeq ($(strip $(BENCHOUT)),)
pprof:
	go test . -bench . -benchmem -benchtime $(BENCHTIME) \
		-o skyline.test \
		-cpuprofile cpu.pprof \
		-memprofile mem.pprof \
		| tee bench.pprof.out
else
pprof:
	mkdir -p $(BENCHOUT)
	go test . -bench . -benchmem -benchtime $(BENCHTIME) \
		-o $(BENCHOUT)/skyline.test \
		-cpuprofile $(BENCHOUT)/cpu.pprof \
		-memprofile $(BENCHOUT)/mem.pprof \
		| tee $(BENCHOUT)/bench.pprof.out
endif

bin/gen: gen/*.go
	go build -o $@ $^

bin/display: display/*.go
	go build -o $@ $^
