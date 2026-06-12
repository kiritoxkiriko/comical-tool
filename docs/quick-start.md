# Quick Start

## Self-hosted

```bash
docker compose -f deploy/docker-compose.yml up --build
```

Open:

- Web UI: `http://localhost:3000`
- API health: `http://localhost:8080/api/health`

The default self-hosted stack uses SQLite, Redis, and local object storage.

## Local Builds

```bash
make test
make lint
make build
make worker-dry-run
make web-build
```

`make worker-dry-run` validates the Worker bundle without deploying it.

## CLI

```bash
cd cli
go run ./cmd/comical-cli --base-url http://localhost:8080 short create --url https://example.com
```

