all: bin/gen bin/display

test:
	go test -v .

BENCHTIME:=100ms

ifeq ($(strip $(OUT)),)
bench:
	go test . -bench . -benchmem -benchtime $(BENCHTIME) \
		-o skyline.test \
		-cpuprofile cpu.pprof \
		-memprofile mem.pprof \
		| tee bench.out
else
bench:
	mkdir -p $(OUT)
	go test . -bench . -benchmem -benchtime $(BENCHTIME) \
		-o $(OUT)/skyline.test \
		-cpuprofile $(OUT)/cpu.pprof \
		-memprofile $(OUT)/mem.pprof \
		| tee $(OUT)/bench.out
endif

bin/gen: gen/*.go
	go build -o $@ $^

bin/display: display/*.go
	go build -o $@ $^
