# comical-tool

`comical-tool` is a small utility platform for short links, image hosting,
temporary clipboard snippets, and temporary file storage.

## Layout

```text
server/     Go Hertz API service, binary: comical-tool
cli/        Go Cobra CLI, binary: comical-cli
web/        Next.js + Tailwind web UI
worker/     Cloudflare Worker API adapter
docs/       VitePress documentation
deploy/     Docker and local deployment files
migrations/ SQL migrations
```

## Local Development

```bash
cd server && go test ./...
cd cli && go test ./...
make build
cd web && npm install && npm run build
cd worker && npm install && npm run build
docker compose -f deploy/docker-compose.yml config
```

Start self-hosted dependencies and apps:

```bash
docker compose -f deploy/docker-compose.yml up --build
```

Server defaults to SQLite, Redis-backed cache, and local object storage. Redis is
included in `docker-compose.yml` for cache-compatible development. Storage can
use either `storage.driver = "local"` or `storage.driver = "s3"`. The
self-hosted server can open SQLite, PostgreSQL, or MySQL through
`database.driver`.

## CLI

```bash
cd cli
go run ./cmd/comical-cli short create --url https://example.com
go run ./cmd/comical-cli clip put --content hello --link
go run ./cmd/comical-cli file upload --file ./secret.zip --password open --max-visits 3
go run ./cmd/comical-cli file download <asset-id> --password open --output ./download.bin
```

## Short Link Domains

Short links support independent domains mapped to a canonical app path.

```toml
[modules.short_link]
domain_mappings = { "s.tool.sqlboy.me" = "https://tool.sqlboy.me/short" }
```

With this mapping, `https://s.tool.sqlboy.me/abc123` maps to the same
short-link resource as `https://tool.sqlboy.me/short/abc123`.

## Cloudflare

Cloudflare support is split into:

- `worker/`: API adapter deployed with Wrangler, backed by D1, R2, and KV.
- `web/`: Next.js deployed to Workers through `@opennextjs/cloudflare`.

GitHub Actions Cloudflare deployment is intentionally disabled. Validate with
dry-runs and builds by default; only deploy manually when the current task
explicitly asks for deployment.

Manual deployment requires:

- `CLOUDFLARE_API_TOKEN`
- `CLOUDFLARE_ACCOUNT_ID`

Required repository variable:

- `NEXT_PUBLIC_API_BASE_URL` defaults to `https://tool.sqlboy.me`

Create Cloudflare D1, R2, and KV resources before production deployment and
verify the binding IDs in `worker/wrangler.jsonc`.

## Documentation

```bash
cd docs
npm install
npm run dev
```

Docs entry points:

- `docs/guide/index.md` for quick start, APIs, CLI, and migrations.
- `docs/api.md` and `docs/cli.md` for command-level API and CLI usage.
- `docs/configuration.md`, `docs/docker.md`, `docs/migrations.md`, and
  `docs/storage.md` for self-hosted operations.
- `docs/cloudflare.md` for Cloudflare verification, explicit manual deployment,
  and runtime differences.
