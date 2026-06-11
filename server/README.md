# server

Go Hertz API service. The binary is `comical-tool`.

```bash
go test ./...
go run ./cmd/comical-tool -config ../deploy/config.example.toml
```

Reusable pure Go logic belongs in `server/pkg`. Runtime code, HTTP handlers,
repository implementations, and storage adapters belong in `server/internal`.
