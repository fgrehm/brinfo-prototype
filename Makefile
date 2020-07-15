.PHONY: build test

build: brinfo-scrape
.PHONY: build

brinfo-scrape: */**/*.go go.mod go.sum
	go build -o brinfo-scrape ./cmd/brinfo-scrape/...

ci: build test lint
.PHONY: ci

test:
	go test ./...
.PHONY: test

# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0
lint:
	golangci-lint run ./...
