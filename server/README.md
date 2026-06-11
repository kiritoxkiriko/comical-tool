# server

Go Hertz API service. The binary is `comical-tool`.

```bash
go test ./...
go build -o ../bin/comical-tool ./cmd/comical-tool
../bin/comical-tool -config ../deploy/config.example.toml
```

Reusable pure Go logic belongs in `server/pkg`. Runtime code, HTTP handlers,
repository implementations, and storage adapters belong in `server/internal`.
The Hertz HTTP layer follows the `biz` layout: `server/internal/biz/router`
registers routes, `server/internal/biz/handler` contains handlers, and
`server/internal/biz/middleware` contains Hertz middleware.
