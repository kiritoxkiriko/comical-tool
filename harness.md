# comical-tool Harness

This file is the repeatable local harness for validation, release checks, and manual Cloudflare deployment. Keep it in sync with `Makefile`, package scripts, and deployment docs.

## Prerequisites

- Go 1.26.x.
- Node.js and npm.
- Docker for dependency services and compose validation.
- Wrangler authenticated locally for manual Cloudflare deployment.
- Project npm installs use `https://registry.npmjs.org/`; keep `.npmrc` files aligned if package-lock files are regenerated.

## Quick Validation

Run this before a focused code change is considered ready:

```bash
make test
make lint
```

## Full Validation

Run this before push or deployment:

```bash
make test
make lint
make build
make worker-dry-run
make web-build
make docker-config
cd docs && npm run build
```

## Local Runtime

Use Docker for dependencies by default. SQLite mode must still work without a database container.

```bash
docker compose -f deploy/docker-compose.yml up -d
cd server && go test ./...
cd web && npm run dev
```

## Cloudflare Manual Deployment

Deploy the Worker API and Web app with Wrangler-backed package scripts:

```bash
cd worker && npm run deploy
cd web && npm run deploy
```

After deploy, smoke-check:

```bash
curl -I https://tool.sqlboy.me
curl -i https://tool.sqlboy.me/api/health
curl -i https://tool.sqlboy.me/api/admin/cleanup -X POST
```

The cleanup endpoint must reject unauthenticated requests. Set `ADMIN_TOKEN` through Wrangler secrets when admin automation needs to run.
