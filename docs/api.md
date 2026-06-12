# API

Business APIs return:

```json
{ "data": {} }
```

Errors return:

```json
{
  "error": {
    "code": "bad_request",
    "message": "invalid request",
    "request_id": "..."
  }
}
```

Health checks and binary asset downloads are not wrapped.

## Short Links

```bash
curl -sS http://localhost:8080/api/short-links \
  -H 'content-type: application/json' \
  -d '{"target_url":"https://example.com","custom_slug":"demo","ttl":"168h"}'
```

Routes:

- `POST /api/short-links`
- `POST /api/short-links/{slug}/revoke`
- `GET /short/{slug}`
- `GET /{slug}` for short-domain mappings

Successful redirects are recorded in `access_events` as `short_link` /
`redirect` events.

## Images

```bash
curl -sS http://localhost:8080/api/images \
  -F file=@./image.png \
  -F ttl=720h \
  -F link=true
```

Routes:

- `POST /api/images`
- `GET /api/images`
- `GET /api/assets/{id}`
- `DELETE /api/images/{id}`

## Clipboard

```bash
curl -sS http://localhost:8080/api/clip \
  -H 'content-type: application/json' \
  -d '{"content":"hello","password":"open","max_visits":5,"ttl":"1h","link":true}'
```

Read with:

```bash
curl -sS 'http://localhost:8080/api/clip/{id}?password=open'
```

Routes:

- `POST /api/clip`
- `GET /api/clip/{id}`
- `DELETE /api/clip/{id}`

## Files

```bash
curl -sS http://localhost:8080/api/files \
  -F file=@./archive.zip \
  -F ttl=168h \
  -F password=open \
  -F max_visits=3 \
  -F link=true
```

Download with:

```bash
curl -LO 'http://localhost:8080/api/assets/{id}?password=open'
```

Routes:

- `POST /api/files`
- `GET /api/files`
- `GET /api/assets/{id}?password=...`
- `DELETE /api/files/{id}`

## Admin

```bash
curl -sS -X POST http://localhost:8080/api/admin/cleanup \
  -H 'authorization: Bearer <admin-token>'
```
