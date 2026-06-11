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
