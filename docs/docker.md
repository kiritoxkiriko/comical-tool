# Docker

Validate the compose file:

```bash
docker compose -f deploy/docker-compose.yml config
```

Start the self-hosted stack:

```bash
docker compose -f deploy/docker-compose.yml up --build
```

Services:

- `server` on port `8080`
- `web` on port `3000`
- `redis` on port `6379`

Optional database profiles:

```bash
docker compose -f deploy/docker-compose.yml --profile postgres up -d postgres
docker compose -f deploy/docker-compose.yml --profile mysql up -d mysql
```

The default volume `comical-data` stores the SQLite database and local objects.

For production self-hosting:

- prefer PostgreSQL or MySQL for multi-instance deployments.
- use S3-compatible storage when objects need to outlive a single host.
- back up the database and object storage together.
- place a reverse proxy in front of `server` and `web`.

