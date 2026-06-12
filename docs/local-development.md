# Local Development

## Dependencies

Use Docker for dependency services:

```bash
docker compose -f deploy/docker-compose.yml up -d redis
```

SQLite mode does not require a database container.

## Server

```bash
cd server
GOTOOLCHAIN=go1.26.0 go test ./...
GOTOOLCHAIN=go1.26.0 go run ./cmd/comical-tool -config ../deploy/config.example.toml
```

## Web

```bash
cd web
npm install
npm run dev
```

## CLI

```bash
cd cli
GOTOOLCHAIN=go1.26.0 go test ./...
GOTOOLCHAIN=go1.26.0 go run ./cmd/comical-cli --help
```

## Worker

```bash
cd worker
npm install
npm run cf-typegen
npm run build
npm run dry-run
```

Do not run deployment commands unless the current task explicitly requests a
manual deployment.

