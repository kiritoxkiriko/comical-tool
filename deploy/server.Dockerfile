FROM golang:1.26-bookworm AS build

WORKDIR /src
COPY go.work ./
COPY server/go.mod server/go.sum ./server/
RUN cd server && GOWORK=off go mod download
COPY server ./server
RUN cd server && GOWORK=off CGO_ENABLED=0 go build -o /out/comical-tool ./cmd/comical-tool

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app
COPY --from=build /out/comical-tool /app/comical-tool
COPY deploy/config.example.toml /app/config.toml
EXPOSE 8080
ENTRYPOINT ["/app/comical-tool", "-config", "/app/config.toml"]
