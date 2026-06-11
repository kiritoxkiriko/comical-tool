package job

import (
	"context"
	"log"
	"time"

	"github.com/kiritoxkiriko/comical-tool/server/internal/repository"
)

type Cleaner interface {
	CleanupExpired(context.Context) (repository.CleanupResult, error)
}

func StartCleanup(ctx context.Context, interval time.Duration, cleaner Cleaner, logger *log.Logger) {
	if interval <= 0 {
		return
	}
	go func() {
		runCleanup(ctx, cleaner, logger)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runCleanup(ctx, cleaner, logger)
			}
		}
	}()
}

func runCleanup(ctx context.Context, cleaner Cleaner, logger *log.Logger) {
	result, err := cleaner.CleanupExpired(ctx)
	if err != nil {
		logger.Printf("cleanup failed: %v", err)
		return
	}
	logger.Printf("cleanup completed: assets=%d clipboard=%d short_links=%d", result.Assets, result.Clipboard, result.ShortLinks)
}
