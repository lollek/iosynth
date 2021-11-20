all: build
.PHONY: all

run:	build
	./iosynth
.PHONY: run

build:
	go build
.PHONY: build

test:
	go test
.PHONY: test

format:
	go fmt
.PHONT: format
