bin/gen: gen/*.go
	go build -o $@ $^
