.PHONY: build test lint run docker-build

build:
	go build -o bin/event-timeline ./cmd/server

test:
	go test ./... -v -race -cover

lint:
	golangci-lint run ./...

run:
	go run ./cmd/server

docker-build:
	docker build -t aevum/event-timeline:latest .
