package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemorySetGetDelete(t *testing.T) {
	store := NewMemory()
	ctx := context.Background()
	if err := store.Set(ctx, "key", []byte("value"), 0); err != nil {
		t.Fatal(err)
	}
	got, ok, err := store.Get(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || string(got) != "value" {
		t.Fatalf("expected cached value, got ok=%v value=%q", ok, string(got))
	}
	if err := store.Delete(ctx, "key"); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := store.Get(ctx, "key"); err != nil || ok {
		t.Fatalf("expected deleted cache miss, ok=%v err=%v", ok, err)
	}
}

func TestMemoryHonorsTTL(t *testing.T) {
	store := NewMemory()
	ctx := context.Background()
	if err := store.Set(ctx, "key", []byte("value"), time.Nanosecond); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond)
	if _, ok, err := store.Get(ctx, "key"); err != nil || ok {
		t.Fatalf("expected expired cache miss, ok=%v err=%v", ok, err)
	}
}
