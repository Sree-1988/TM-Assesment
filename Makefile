.PHONY: help build run test lint clean fmt

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "Service Registry - Available Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'

## build: Build the server binary
build:
	go build -o bin/server ./cmd/server

## run: Run the server locally
run:
	go run ./cmd/server

## test: Run all tests
test:
	go test -v ./...

## lint: Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

## fmt: Format code
fmt:
	go fmt ./...

## clean: Remove build artifacts
clean:
	rm -rf bin/

## tidy: Tidy go modules
tidy:
	go mod tidy
