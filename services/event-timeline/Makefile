.PHONY: build test test-unit test-integration lint run docker-build

build:
	go build -o bin/event-timeline ./cmd/server

test:
	$(MAKE) test-unit

test-unit:
	go test $(shell go list ./... | grep -v '/tests/integration') -v -race -cover

test-integration:
	go test ./tests/integration -tags=integration -v -race -cover

lint:
	golangci-lint run ./...

run:
	go run ./cmd/server

docker-build:
	docker build -t aevum/event-timeline:latest .
