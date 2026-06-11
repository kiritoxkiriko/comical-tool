# Guide

## Self-hosted

```bash
docker compose -f deploy/docker-compose.yml up --build
```

## API

Business APIs return a stable JSON envelope:

```json
{
  "data": {}
}
```

Error responses use:

```json
{
  "error": {
    "code": "bad_request",
    "message": "invalid request",
    "request_id": "..."
  }
}
```

Health checks and binary asset downloads are intentionally not wrapped.

- `POST /api/short-links`
- `POST /api/short-links/:slug/revoke`
- `POST /api/images`
- `GET /api/images`
- `POST /api/clip`
- `GET /api/clip/:id`
- `POST /api/files`
- `GET /api/files`

Expired or revoked resources return `410 Gone`. Invalid input, including uploads
larger than the configured module limit, returns `400 Bad Request`.

## Configuration

Self-hosted config is TOML. The server applies `server.max_body_bytes` as the
global request body limit, then enforces module limits:

- `modules.image_hosting.max_bytes` for `POST /api/images`.
- `modules.file_stash.max_bytes` for `POST /api/files`.

See `deploy/config.example.toml` for the default values.

Administrative HTTP APIs require `Authorization: Bearer <security.admin_token>`.
For CLI usage, pass `--token` or set `COMICAL_ADMIN_TOKEN`.

Database config uses `[database]`. The self-hosted server can open and migrate
SQLite, PostgreSQL, and MySQL:

- `driver = "sqlite"` with `dsn = "file:/data/comical.db?_foreign_keys=on"`.
- `driver = "postgres"` with a PostgreSQL URL.
- `driver = "mysql"` with a MySQL DSN; include `parseTime=true`.

Self-hosted cleanup runs automatically when `cleanup.enabled = true`. The
default interval is `30m`; set `cleanup.interval` in TOML to change it.

Cache config uses `[cache]`. The self-hosted server opens and pings the cache at
startup:

- `driver = "redis"` with `dsn = "redis://redis:6379/0"` for Docker/local Redis.
- `driver = "memory"` for process-local cache during isolated development.

Storage config uses `[storage]`:

- `driver = "local"` with `local_dir` for the default filesystem object store.
- `driver = "s3"` with `s3_endpoint`, `s3_region`, `s3_bucket`,
  `s3_access_key_id`, `s3_secret_access_key`, and optional
  `s3_use_path_style = true` for S3-compatible providers.

## CLI

```bash
comical-cli short create --url https://example.com
comical-cli image upload --file ./image.png --link
comical-cli clip put --content "hello" --ttl 1h
comical-cli file upload --file ./archive.zip --ttl 168h
comical-cli file download <asset-id> --output ./archive.zip
```

## Cloudflare

Cloudflare deployment support is available, but GitHub Actions deployment is
disabled. Use dry-runs and builds for validation unless the current task
explicitly asks for a manual deployment.

Before production deployment, create D1, R2, and KV resources and update
`worker/wrangler.jsonc`.

See [Cloudflare](../cloudflare.md) for Cloudflare verification, explicit manual
deployment, runtime differences, and smoke checks.

## Migrations

Migration directories:

- `migrations/sqlite` for local default self-hosting.
- `migrations/postgres` for PostgreSQL self-hosting.
- `migrations/mysql` for MySQL self-hosting.
- `migrations/d1` for Cloudflare Worker metadata.

PostgreSQL/MySQL smoke tests are opt-in:

```bash
docker compose -f deploy/docker-compose.yml --profile postgres --profile mysql up -d postgres mysql
COMICAL_TEST_POSTGRES_DSN='postgres://comical:comical@127.0.0.1:15432/comical?sslmode=disable' \
COMICAL_TEST_MYSQL_DSN='comical:comical@tcp(127.0.0.1:13306)/comical?parseTime=true' \
  go test ./server/internal/repository
```
