# Troubleshooting

## Worker Build Fails After Type Generation

Run:

```bash
cd worker
npm run cf-typegen
npm run build
```

The project relies on Wrangler-generated `*.wasm` module declarations. Avoid
adding another generic `declare module "*.wasm"` declaration.

## Go Uses The Wrong Toolchain

Use:

```bash
GOTOOLCHAIN=go1.26.0 go test ./...
```

The root `Makefile` already applies this setting for Go commands.

## File Download Returns 403

The file may have been uploaded with a password. Pass it as a query parameter or
CLI flag:

```bash
curl -LO 'http://localhost:8080/api/assets/{id}?password=open'
comical-cli file download <asset-id> --password open --output ./file.bin
```

## File Or Clipboard Returns 410

The resource is expired, revoked, deleted, or has reached its visit limit.

## Image Upload Returns 400

The image hosting endpoint only accepts multipart files with an `image/*` MIME
type and enforces `modules.image_hosting.max_bytes`.

## Cloudflare Deployment

GitHub Actions deployment is disabled. Do not run `wrangler deploy` unless the
current task explicitly asks for a manual deployment. Use:

```bash
make worker-dry-run
make web-build
```
