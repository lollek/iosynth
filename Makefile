all: build
.PHONY: all

build:
	go build
.PHONY: build

test:
	go test
.PHONY: test

format:
	go fmt
.PHONT: format
