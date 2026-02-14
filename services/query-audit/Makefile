.PHONY: help build run test test-unit test-integration lint fmt docker-build docker-run clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/server ./cmd/server

run: ## Run the application
	go run ./cmd/server

test: ## Run all tests
	go test ./... -v -race

test-unit: ## Run unit tests
	go test ./tests/unit/... -v

test-integration: ## Run integration tests
	go test ./tests/integration/... -v

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	goimports -w .

docker-build: ## Build Docker image
	docker build -t query-audit:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 \
		-e ELASTICSEARCH_URLS=http://elasticsearch:9200 \
		-e EVENT_TIMELINE_URL=http://event-timeline:8081 \
		-e DECISION_ENGINE_URL=http://decision-engine:5000 \
		query-audit:latest

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

deps: ## Download dependencies
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
