package repository

import (
	"context"
	"testing"
	"time"

	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

func TestSQLiteShortLink(t *testing.T) {
	repo := openTestRepo(t)
	link := domain.ShortLink{
		ID:        "id-1",
		Slug:      "abc123",
		TargetURL: "https://example.com",
	}
	if err := repo.CreateShortLink(context.Background(), link); err != nil {
		t.Fatal(err)
	}
	got, err := repo.FindShortLink(context.Background(), "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if got.TargetURL != link.TargetURL {
		t.Fatalf("expected %q, got %q", link.TargetURL, got.TargetURL)
	}
}

func TestSQLiteClipboard(t *testing.T) {
	repo := openTestRepo(t)
	expires := time.Now().UTC().Add(time.Hour)
	item := domain.ClipboardItem{ID: "clip-1", Content: "hello", ExpiresAt: &expires}
	if err := repo.CreateClipboard(context.Background(), item); err != nil {
		t.Fatal(err)
	}
	if err := repo.IncrementClipboardVisit(context.Background(), item.ID); err != nil {
		t.Fatal(err)
	}
	got, err := repo.FindClipboard(context.Background(), item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.VisitCount != 1 {
		t.Fatalf("expected 1 visit, got %d", got.VisitCount)
	}
}

func openTestRepo(t *testing.T) *SQLite {
	t.Helper()
	repo, err := OpenSQLite(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	if err := repo.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	return repo
}
