.PHONY: help build test clean restore run docker-build docker-run lint

SOLUTION := DecisionEngine.sln
PROJECT := src/Aevum.DecisionEngine.Api/Aevum.DecisionEngine.Api.csproj
DOCKER_IMAGE := aevum/decision-engine
DOCKER_TAG := latest

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

restore: ## Restore NuGet packages
	dotnet restore $(SOLUTION)

build: restore ## Build the solution
	dotnet build $(SOLUTION) -c Release --no-restore

test: build ## Run all tests
	dotnet test $(SOLUTION) -c Release --no-build --verbosity normal

test-unit: build ## Run unit tests only
	dotnet test tests/Aevum.DecisionEngine.Domain.Tests/Aevum.DecisionEngine.Domain.Tests.csproj -c Release --no-build
	dotnet test tests/Aevum.DecisionEngine.Application.Tests/Aevum.DecisionEngine.Application.Tests.csproj -c Release --no-build

test-integration: build ## Run integration tests only
	dotnet test tests/Aevum.DecisionEngine.Integration.Tests/Aevum.DecisionEngine.Integration.Tests.csproj -c Release --no-build

test-coverage: ## Run tests with coverage
	dotnet test $(SOLUTION) -c Release --collect:"XPlat Code Coverage" --results-directory ./coverage

clean: ## Clean build artifacts
	dotnet clean $(SOLUTION)
	rm -rf */bin */obj coverage

run: ## Run the API locally
	dotnet run --project $(PROJECT)

watch: ## Run with hot reload
	dotnet watch --project $(PROJECT)

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	docker run --rm -it -p 8080:8080 -p 9090:9090 \
		-e MongoDB__ConnectionString=mongodb://host.docker.internal:27017 \
		-e MongoDB__DatabaseName=decision_engine \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

lint: ## Run code analysis
	dotnet format $(SOLUTION) --verify-no-changes

format: ## Format code
	dotnet format $(SOLUTION)

deps-update: ## Update all NuGet packages
	dotnet list $(SOLUTION) package --outdated
