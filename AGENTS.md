# comical-tool Agent Boundaries

This file defines project boundaries and agent behavior only. It is not the development plan. The development plan lives in `docs/development-plan.md`.

## Source Of Truth

- Development plan: `docs/development-plan.md`
- Repository structure, naming, and module boundaries must stay aligned with the development plan.
- If implementation reveals that the plan is no longer valid, update the development plan before changing code.

## Repository Boundaries

The repository is organized by product/runtime domains at the top level. Each domain then follows the standard structure for its implementation language.

Allowed top-level directories and files:

```text
server/         # Go Hertz API service
cli/            # Go CLI
web/            # Next.js frontend
worker/         # Cloudflare Worker adapter
docs/           # Project docs and VitePress docs site
deploy/         # Docker, docker-compose, wrangler, reverse proxy config
migrations/     # sqlite/postgres/mysql migrations
scripts/        # Local development, migration, build, and cleanup scripts
.npmrc          # Root npm registry policy for local and CI installs
harness.md      # Repeatable validation and deployment harness
go.work         # Optional workspace for Go modules under server/ and cli/
```

Do not:

- Do not add `apps/`, `apps/web`, or `apps/api`.
- Do not add a top-level `core/`, `cmd/`, `internal/`, or `pkg/` directory.
- Do not import `server/internal` from `cli`, `worker`, or any external module.
- Do not put reusable domain/policy logic under `server/internal` if it must also be used by CLI or WASM builds.
- Do not let business logic depend directly on concrete database, cache, or object storage SDKs.

## Domain Structure

- `server/` is a Go module. Use `server/cmd/comical-tool`, `server/internal`, and `server/pkg`.
- Server HTTP code follows Hertz-style layout under `server/internal/biz`: handlers in `biz/handler`, route registration in `biz/router`, and middleware in `biz/middleware`.
- `cli/` is a Go module. Use `cli/cmd/comical-cli` and `cli/internal`; build the command tree with Cobra.
- `web/` is a Next.js project and follows normal Next.js/Tailwind conventions.
- `worker/` is a Cloudflare Worker project and follows normal Workers conventions.
- Cross-module Go development can use root `go.work`, but domain code stays inside its domain directory.

## Runtime Boundaries

- `server/cmd/comical-tool` starts the self-hosted Go Hertz backend.
- `web` is the main frontend and uses Next.js.
- `cli/cmd/comical-cli` is the Cobra-based CLI and talks to the service through HTTP APIs.
- `worker` is a Cloudflare-specific adapter and must not run Hertz directly.
- Shared pure Go logic used by the server belongs in `server/pkg`, for example `server/pkg/domain`, `server/pkg/policy`, and `server/pkg/apperror`.
- `server/pkg` must not depend on Hertz, sqlx, Redis, S3, R2, D1, KV, or other runtime/infrastructure clients.
- Worker-side WASM can reuse selected `server/pkg` logic, but DB, storage, cache, and HTTP runtime code stay outside WASM.

## Config And Secrets

- Configuration uses TOML and can be overridden by environment variables.
- Configuration structs should be reusable by server and CLI without importing `server/internal`.
- Do not commit secrets, tokens, bucket secrets, database passwords, or real production DSNs.
- When adding config keys, update example config and docs in the same change.

## Cloudflare Boundaries

This project has project-scoped Cloudflare skills installed under `.agents/skills/`. For Cloudflare Worker, D1, R2, KV, OpenNext, or wrangler work, read the relevant skills first.

Keep Cloudflare setup project-scoped:

- Do not install Cloudflare skills globally unless the user explicitly asks.
- Do not run `codex mcp add ...` without confirmation; it writes to the user-level Codex config.
- Use `.codex/cloudflare-mcp.toml` as this project's local MCP config snippet.
- If Cloudflare MCP login is needed, confirm that the user wants to use Cloudflare MCP for this project first.

## Development Rules

- Before changing code, read `docs/development-plan.md` and the README in the target domain directory.
- For validation or deployment work, read `harness.md` and keep it aligned with Makefile/package script changes.
- For shared schema, migrations, config structures, deployment files, or cross-runtime interfaces, check the impact area first.
- For new features, add tests for `server/pkg` and `server/internal` first, then wire HTTP, Web, CLI, or Worker behavior.
- Local development should use Docker for dependencies by default; SQLite mode must not require a database container.
- Do not delete existing user files or untracked content unless the user explicitly asks.

## Verification

Choose verification based on the changed area:

- Harness: follow `harness.md` for the repeatable command sequence.
- Server: `cd server && go test ./...`
- CLI: `cd cli && go test ./...`
- Web: `cd web && npm run lint`, `npm run test`, `npm run build`
- Docs: VitePress build
- Docker: `docker compose config`, and `docker compose up` when needed
- Cloudflare: Wrangler/OpenNext build or dry-run

If verification cannot run, state the reason. Do not claim it passed.
