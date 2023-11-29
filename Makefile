.PHONY: build
build:
				go build -v ./cmd/app

.PHONY: test
test:
				go test -v -race ./...

.PHONY: linter
linter:
				golangci-lint run ./... --config=./.golangci.yml

.PHONY: up
up:
				.\app

.PHONY: test-up
test-up:
				.\app -test-server=true

.DEFAULT_GOAL := build