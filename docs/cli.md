# CLI

The CLI binary is `comical-cli`.

```bash
cd cli
go build -o ../bin/comical-cli ./cmd/comical-cli
```

Global flags:

- `--base-url`
- `--token`
- `--output`, `-o`

## Short Links

```bash
comical-cli short create --url https://example.com --slug demo --ttl 168h
comical-cli short revoke demo
```

## Images

```bash
comical-cli image upload --file ./image.png --ttl 720h --link
comical-cli image list
comical-cli image delete <asset-id>
```

Image uploads set the multipart file `Content-Type` from the local file
extension, so common formats such as `.png`, `.jpg`, `.gif`, and `.webp` satisfy
the image hosting MIME validation.

## Clipboard

```bash
comical-cli clip put --content hello --password open --max-visits 5 --ttl 1h --link
comical-cli clip get <id> --password open
comical-cli clip delete <id>
```

## Files

```bash
comical-cli file upload --file ./archive.zip --ttl 168h --password open --max-visits 3 --link
comical-cli file list
comical-cli file download <asset-id> --password open --output ./archive.zip
comical-cli file delete <asset-id>
```

## Admin

```bash
comical-cli --token <admin-token> admin cleanup
```
