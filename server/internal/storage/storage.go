package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
)

type ObjectInfo struct {
	Size        int64
	ContentType string
}

type Store interface {
	Put(ctx context.Context, key string, body io.Reader) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	Head(ctx context.Context, key string) (ObjectInfo, error)
}

func Open(ctx context.Context, cfg config.StorageConfig) (Store, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Driver)) {
	case "", "local":
		store := NewLocal(cfg.LocalDir)
		return store, store.Ensure(ctx)
	case "s3":
		return NewS3(ctx, cfg)
	default:
		return nil, fmt.Errorf("unsupported storage driver %q", cfg.Driver)
	}
}
