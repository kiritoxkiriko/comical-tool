# cli

Go Cobra CLI. The binary is `comical-cli`.

```bash
go test ./...
go build -o ../bin/comical-cli ./cmd/comical-cli
../bin/comical-cli --help
```

The CLI talks to the server through HTTP APIs and must not import
`server/internal`.
