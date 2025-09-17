.PHONY: build run test clean docker-build docker-run docker-stop

# Build the application
build:
	go build -o bin/comical-tool .

# Run the application
run:
	go run .

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
