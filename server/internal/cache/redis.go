package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func openRedis(ctx context.Context, dsn string) (*Redis, error) {
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	store := &Redis{client: redis.NewClient(opts)}
	if err := store.Ping(ctx); err != nil {
		_ = store.Close()
		return nil, err
	}
	return store, nil
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, bool, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err == nil {
		return data, true, nil
	}
	if err == redis.Nil {
		return nil, false, nil
	}
	return nil, false, fmt.Errorf("redis get %q: %w", key, err)
}

func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) Close() error {
	return r.client.Close()
}
