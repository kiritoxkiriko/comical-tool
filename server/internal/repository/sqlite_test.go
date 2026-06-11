package repository

import (
	"context"
	"fmt"
	"os"
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
	testClipboard(t, repo, "sqlite")
}

func TestExternalDatabases(t *testing.T) {
	tests := []struct {
		name   string
		driver string
		env    string
	}{
		{name: "postgres", driver: "postgres", env: "COMICAL_TEST_POSTGRES_DSN"},
		{name: "mysql", driver: "mysql", env: "COMICAL_TEST_MYSQL_DSN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := os.Getenv(tt.env)
			if dsn == "" {
				t.Skipf("%s is not set", tt.env)
			}
			repo, err := Open(tt.driver, dsn)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = repo.Close() })
			if err := repo.Migrate(context.Background()); err != nil {
				t.Fatal(err)
			}
			suffix := fmt.Sprintf("%s-%d", tt.name, time.Now().UnixNano())
			testShortLink(t, repo, suffix)
			testClipboard(t, repo, suffix)
		})
	}
}

func testShortLink(t *testing.T, repo *Store, suffix string) {
	t.Helper()
	link := domain.ShortLink{
		ID:        fmt.Sprintf("id-%s", suffix),
		Slug:      fmt.Sprintf("slug-%s", suffix),
		TargetURL: "https://example.com",
	}
	if err := repo.CreateShortLink(context.Background(), link); err != nil {
		t.Fatal(err)
	}
	got, err := repo.FindShortLink(context.Background(), link.Slug)
	if err != nil {
		t.Fatal(err)
	}
	if got.TargetURL != link.TargetURL {
		t.Fatalf("expected %q, got %q", link.TargetURL, got.TargetURL)
	}
}

func testClipboard(t *testing.T, repo *Store, suffix string) {
	t.Helper()
	expires := time.Now().UTC().Add(time.Hour)
	item := domain.ClipboardItem{ID: "clip-" + suffix, Content: "hello", ExpiresAt: &expires}
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

func TestNormalizeDriver(t *testing.T) {
	tests := map[string]string{
		"":           "sqlite",
		"sqlite":     "sqlite",
		"sqlite3":    "sqlite",
		"postgres":   "pgx",
		"postgresql": "pgx",
		"pgx":        "pgx",
		"mysql":      "mysql",
	}
	for input, want := range tests {
		got, err := normalizeDriver(input)
		if err != nil {
			t.Fatalf("normalizeDriver(%q): %v", input, err)
		}
		if got != want {
			t.Fatalf("normalizeDriver(%q) = %q, want %q", input, got, want)
		}
	}
	if _, err := normalizeDriver("oracle"); err == nil {
		t.Fatal("expected unsupported driver error")
	}
}

func TestDBTimeScansSupportedValues(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Microsecond)
	inputs := []any{
		now,
		now.Format(time.RFC3339Nano),
		[]byte(now.Format("2006-01-02 15:04:05.999999")),
	}
	for _, input := range inputs {
		var scanned dbTime
		if err := scanned.Scan(input); err != nil {
			t.Fatalf("Scan(%T): %v", input, err)
		}
		if !scanned.Valid {
			t.Fatalf("Scan(%T) returned invalid time", input)
		}
	}
}

func openTestRepo(t *testing.T) *Store {
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
