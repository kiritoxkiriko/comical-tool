.PHONY: build run test clean docker-build docker-run docker-stop

# Project configuration
MODULE_NAME ?= github.com/kiritoxkiriko/comical-tool

# Build the application
build:
	@mkdir -p bin
	go build -ldflags "-X $(MODULE_NAME)/cmd.version=$(VERSION) -X $(MODULE_NAME)/cmd.buildTime=$(BUILD_TIME) -X $(MODULE_NAME)/cmd.gitCommit=$(GIT_COMMIT)" -o bin/comical-tool .

# Build with version info
VERSION ?= 1.0.0
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build for production with full git hash
build-prod:
	@mkdir -p bin
	go build -ldflags "-X $(MODULE_NAME)/cmd.version=$(VERSION) -X $(MODULE_NAME)/cmd.buildTime=$(BUILD_TIME) -X $(MODULE_NAME)/cmd.gitCommit=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")" -o bin/comical-tool .

# Build for development
build-dev:
	@mkdir -p bin
	go build -ldflags "-X $(MODULE_NAME)/cmd.version=$(VERSION)-dev -X $(MODULE_NAME)/cmd.buildTime=$(BUILD_TIME) -X $(MODULE_NAME)/cmd.gitCommit=$(GIT_COMMIT)" -o bin/comical-tool .

# Run the application
run:
	go run .

# Run the server
run-server:
	go run . server

# Show help
help:
	go run . --help

# Show server help
help-server:
	go run . server --help

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t comical-tool .

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose
docker-stop:
	docker-compose down

# Run with Docker Compose and rebuild
docker-rebuild:
	docker-compose up --build -d

# View logs
logs:
	docker-compose logs -f app

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate IDL files
generate-idl:
	chmod +x scripts/generate_idl.sh
	./scripts/generate_idl.sh

# Generate go.mod
mod-init:
	go mod init github.com/kiritoxkiriko/comical-tool

# Update dependencies
mod-update:
	go get -u ./...
	go mod tidy
