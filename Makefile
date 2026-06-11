GO ?= GOTOOLCHAIN=go1.26.0 go
GOLANGCI_LINT ?= golangci-lint

.PHONY: test lint build docker-config worker-dry-run web-build

test:
	cd server && $(GO) test ./...
	cd cli && $(GO) test ./...
	cd web && npm run test
	cd worker && npm run build

lint:
	cd server && $(GOLANGCI_LINT) run ./...
	cd cli && $(GOLANGCI_LINT) run ./...
	cd web && npm run lint
	cd worker && npm run lint

build:
	mkdir -p bin
	cd server && $(GO) build -o ../bin/comical-tool ./cmd/comical-tool
	cd cli && $(GO) build -o ../bin/comical-cli ./cmd/comical-cli
	cd web && npm run build

docker-config:
	docker compose -f deploy/docker-compose.yml config

worker-dry-run:
	cd worker && npm run dry-run

web-build:
	cd web && npm run cf:build
