package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kiritoxkiriko/comical-tool/server/internal/config"
)

type Store interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}

func Open(ctx context.Context, cfg config.CacheConfig) (Store, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Driver)) {
	case "", "memory", "inmemory":
		store := NewMemory()
		return store, store.Ping(ctx)
	case "redis":
		return openRedis(ctx, cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported cache driver %q", cfg.Driver)
	}
}
