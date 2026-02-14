# Makefile for gocheck

.PHONY: help build run test clean install deps

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

deps: ## Download dependencies
	go mod download
	go mod verify

build: ## Build the binary
	go build -o gocheck.exe -v .

run: ## Run the application with sample data
	go run . testdata/sample.csv

run-verbose: ## Run with verbose output
	go run . -v testdata/sample.csv

run-json: ## Run with JSON output
	go run . -f json testdata/sample.csv

run-products: ## Run with products sample data
	go run . testdata/products.csv

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -f gocheck.exe gocheck
	rm -f coverage.out coverage.html
	go clean

install: build ## Install the binary to GOPATH/bin
	go install .

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

all: clean deps fmt vet build test ## Run all checks and build
