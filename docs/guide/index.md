# Guide

## Self-hosted

```bash
docker compose -f deploy/docker-compose.yml up --build
```

## API

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

## CLI

```bash
comical-cli short create --url https://example.com
comical-cli image upload --file ./image.png --link
comical-cli clip put --content "hello" --ttl 1h
comical-cli file upload --file ./archive.zip --ttl 168h
```

## Cloudflare

Deploy the API Worker from `worker/` and the Web UI from `web/`.

Before production deployment, create D1, R2, and KV resources and update
`worker/wrangler.jsonc`.

See [Cloudflare Deployment](../cloudflare.md) for manual Wrangler deployment,
runtime differences, and smoke checks.

## Migrations

Migration directories:

- `migrations/sqlite` for local default self-hosting.
- `migrations/postgres` for PostgreSQL self-hosting.
- `migrations/mysql` for MySQL self-hosting.
- `migrations/d1` for Cloudflare Worker metadata.
