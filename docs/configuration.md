# Configuration

Self-hosted configuration is TOML. Use `deploy/config.example.toml` as the
starting point.

## Server

```toml
[server]
addr = "0.0.0.0:8080"
public_base_url = "https://tool.sqlboy.me"
max_body_bytes = 104857600
```

`public_base_url` is used when the server generates public short URLs.

## Database

```toml
[database]
driver = "sqlite"
dsn = "file:/data/comical.db?_foreign_keys=on"
```

Supported drivers:

- `sqlite`
- `postgres`
- `mysql`

## Cache

```toml
[cache]
driver = "redis"
dsn = "redis://redis:6379/0"
```

Use `driver = "memory"` for isolated local development without Redis.

## Storage

```toml
[storage]
driver = "local"
local_dir = "/data/objects"
```

Use `driver = "s3"` with the S3 fields from `deploy/config.example.toml` for
S3-compatible object storage.

## Modules

```toml
[modules.short_link]
default_ttl = "168h"
allow_custom_slug = true
domain_mappings = { "s.tool.sqlboy.me" = "https://tool.sqlboy.me/short" }

[modules.image_hosting]
default_ttl = "720h"
max_bytes = 10485760

[modules.clipboard]
default_ttl = "1h"
max_visits = 5

[modules.file_stash]
default_ttl = "168h"
max_bytes = 104857600
```

File uploads can also pass per-resource `password` and `max_visits` values.

