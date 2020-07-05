.PHONY: build test

build: brinfo-scrape

brinfo-scrape: **/*.go go.mod go.sum
	go build -o brinfo-scrape cmd/brinfo-scrape/...

test:
	go test ./...
