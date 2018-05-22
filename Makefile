all: bin/gen bin/display

bin/gen: gen/*.go
	go build -o $@ $^

bin/display: display/*.go
	go build -o $@ $^
