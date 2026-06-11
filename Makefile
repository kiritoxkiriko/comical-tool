.PHONY: test build docker-config worker-dry-run web-build

test:
	cd server && go test ./...
	cd cli && go test ./...
	cd web && npm run test
	cd worker && npm run build

build:
	cd server && go build ./cmd/comical-tool
	cd cli && go build ./cmd/comical-cli
	cd web && npm run build

docker-config:
	docker compose -f deploy/docker-compose.yml config

worker-dry-run:
	cd worker && npm run dry-run

web-build:
	cd web && npm run cf:build
