# cli

Go Cobra CLI. The binary is `comical-cli`.

```bash
go test ./...
go run ./cmd/comical-cli --help
```

The CLI talks to the server through HTTP APIs and must not import
`server/internal`.
