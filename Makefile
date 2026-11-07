.PHONY: run build test lint clean

BINARY=bin/server

run:
	go run ./cmd/server

build:
	go build -o $(BINARY) ./cmd/server

test:
	go test -v -race ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/
