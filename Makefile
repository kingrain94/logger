.PHONY: help test test-coverage test-race lint fmt vet build clean examples deps

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detection
	go test -race ./...

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	go vet ./...

build: ## Build the project
	go build ./...

clean: ## Clean build artifacts and temporary files
	go clean ./...
	rm -f coverage.out coverage.html
	rm -f *.log
	find . -name "*.test" -delete

# Example targets
examples: ## Run all examples
	@echo "Running basic example..."
	cd examples/basic && go run main.go
	@echo "\nRunning configuration example..."
	cd examples/configuration && go run main.go
	@echo "\nRunning advanced example..."
	cd examples/advanced && go run main.go

example-basic: ## Run basic example
	cd examples/basic && go run main.go

example-config: ## Run configuration example
	cd examples/configuration && go run main.go

example-advanced: ## Run advanced example
	cd examples/advanced && go run main.go

example-web: ## Run web server example (runs in background)
	@echo "Starting web server on :8080..."
	@echo "Visit http://localhost:8080, http://localhost:8080/health, http://localhost:8080/error"
	@echo "Press Ctrl+C to stop"
	cd examples/web-server && go run main.go

# Dependency management
deps: ## Download and verify dependencies
	go mod download
	go mod verify

deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

# Quality checks
check: fmt vet lint test ## Run all quality checks

# Benchmark
bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

# Install tools
install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Release preparation
pre-commit: clean fmt vet lint test ## Run all checks before committing

# Docker targets (if you want to add Docker support later)
docker-build: ## Build Docker image
	docker build -t kingrain94/logger .

docker-run: ## Run Docker container
	docker run --rm -it kingrain94/logger
