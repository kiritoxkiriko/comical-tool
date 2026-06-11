package storage

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestLocalPutOpenHeadDelete(t *testing.T) {
	store := NewLocal(t.TempDir())
	ctx := context.Background()
	if err := store.Ensure(ctx); err != nil {
		t.Fatal(err)
	}
	if err := store.Put(ctx, "file/id", bytes.NewBufferString("hello")); err != nil {
		t.Fatal(err)
	}
	info, err := store.Head(ctx, "file/id")
	if err != nil {
		t.Fatal(err)
	}
	if info.Size != 5 {
		t.Fatalf("expected size 5, got %d", info.Size)
	}
	body, err := store.Open(ctx, "file/id")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = body.Close()
	}()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("expected stored body, got %q", string(data))
	}
	if err := store.Delete(ctx, "file/id"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Head(ctx, "file/id"); err == nil {
		t.Fatal("expected deleted object to be missing")
	}
}

func TestLocalRejectsUnsafeKeys(t *testing.T) {
	store := NewLocal(t.TempDir())
	ctx := context.Background()

	unsafeKeys := []string{"../secret", "/tmp/secret", ".."}
	for _, key := range unsafeKeys {
		if err := store.Put(ctx, key, bytes.NewBufferString("bad")); err == nil {
			t.Fatalf("expected key %q to be rejected", key)
		}
	}
}
