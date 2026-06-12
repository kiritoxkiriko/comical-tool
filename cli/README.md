# cli

Go Cobra CLI. The binary is `comical-cli`.

```bash
go test ./...
go build -o ../bin/comical-cli ./cmd/comical-cli
../bin/comical-cli --help
../bin/comical-cli file upload --file ./archive.zip --password open --max-visits 3
../bin/comical-cli file download <asset-id> --password open --output ./archive.zip
```

The CLI talks to the server through HTTP APIs and must not import
`server/internal`.
