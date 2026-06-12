# Migrations

Migration directories:

- `migrations/sqlite`
- `migrations/postgres`
- `migrations/mysql`
- `migrations/d1`

The self-hosted server can run the equivalent schema through
`server/internal/repository`.

The schema includes `access_events` for resource access audit records. Short
link redirects write a `short_link` / `redirect` event in both self-hosted and
Cloudflare runtimes.

The schema also includes `resource_links` so short links generated for images,
clipboard entries, and temporary files can be traced back to their source
resource.

SQLite is the default local mode:

```toml
[database]
driver = "sqlite"
dsn = "file:/data/comical.db?_foreign_keys=on"
```

PostgreSQL smoke test:

```bash
docker compose -f deploy/docker-compose.yml --profile postgres up -d postgres
COMICAL_TEST_POSTGRES_DSN='postgres://comical:comical@127.0.0.1:15432/comical?sslmode=disable' \
  GOTOOLCHAIN=go1.26.0 go test ./server/internal/repository
```

MySQL smoke test:

```bash
docker compose -f deploy/docker-compose.yml --profile mysql up -d mysql
COMICAL_TEST_MYSQL_DSN='comical:comical@tcp(127.0.0.1:13306)/comical?parseTime=true' \
  GOTOOLCHAIN=go1.26.0 go test ./server/internal/repository
```

Cloudflare D1 uses `migrations/d1`.

Local D1 migration smoke test:

```bash
scripts/check-d1-migrations.sh
```
