all: bin/gen bin/display

test:
	go test -v .

bin/gen: gen/*.go
	go build -o $@ $^

bin/display: display/*.go
	go build -o $@ $^
