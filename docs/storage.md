# Storage Backends

The self-hosted server uses `server/internal/storage.Store` behind the service
layer. Business logic does not depend on a concrete provider client.

## Local

```toml
[storage]
driver = "local"
local_dir = "/data/objects"
```

Local storage is useful for SQLite single-host deployments and development.

## S3 Compatible

```toml
[storage]
driver = "s3"
s3_endpoint = "https://s3.example.com"
s3_region = "auto"
s3_bucket = "comical-tool"
s3_access_key_id = "change-me"
s3_secret_access_key = "change-me"
s3_use_path_style = true
```

Use S3-compatible storage for production self-hosting when files need durable
object storage outside the server filesystem.

## Cloudflare R2

Cloudflare Worker deployments use the R2 binding in `worker/wrangler.jsonc`.
The Worker API adapter stores object metadata in D1 and bytes in R2.

