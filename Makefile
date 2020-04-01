.PHONY: build test bench

build:
	go build

test:
	go test

bench:
	go test -bench=.
