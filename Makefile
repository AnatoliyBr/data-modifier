.PHONY: build
build:
				go build -v ./cmd/app

.PHONY: test
test:
				go test -v -race ./...

.PHONY: linter
linter:
				golangci-lint run ./... --config=./.golangci.yml


.DEFAULT_GOAL := build