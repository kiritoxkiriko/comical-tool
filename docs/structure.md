# Project Structure

The repository is organized by runtime domain at the top level. Each domain then
uses its own language conventions.

```text
server/      Go Hertz API service
cli/         Go Cobra CLI
web/         Next.js Web UI
worker/      Cloudflare Worker API adapter
docs/        VitePress documentation
deploy/      Docker and local deployment files
migrations/  SQLite/PostgreSQL/MySQL/D1 migrations
scripts/     Local helper scripts
```

Important boundaries:

- `server/cmd/comical-tool` starts the self-hosted API service.
- `cli/cmd/comical-cli` starts the CLI.
- `server/pkg` contains pure reusable Go logic only.
- `server/internal` contains runtime implementation details and must not be
  imported by `cli` or `worker`.
- `worker` does not run Hertz. It reuses selected `server/pkg/policy` logic
  through Go WASM.

