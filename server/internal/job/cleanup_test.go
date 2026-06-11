package job

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	"github.com/kiritoxkiriko/comical-tool/server/internal/repository"
)

func TestStartCleanupRunsImmediately(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cleaner := &fakeCleaner{called: make(chan struct{}, 1)}
	StartCleanup(ctx, time.Hour, cleaner, log.New(io.Discard, "", 0))
	select {
	case <-cleaner.called:
	case <-time.After(time.Second):
		t.Fatal("expected cleanup to run immediately")
	}
}

type fakeCleaner struct {
	called chan struct{}
}

func (f *fakeCleaner) CleanupExpired(context.Context) (repository.CleanupResult, error) {
	f.called <- struct{}{}
	return repository.CleanupResult{Assets: 1, Clipboard: 2, ShortLinks: 3}, nil
}
