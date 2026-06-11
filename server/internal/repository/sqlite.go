package repository

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var ErrNotFound = errors.New("not found")

type SQLite struct {
	db *sqlx.DB
}

func OpenSQLite(dsn string) (*SQLite, error) {
	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &SQLite{db: db}, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, schema)
	return err
}

const schema = `
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  created_at TEXT NOT NULL
);

INSERT OR IGNORE INTO users (id, created_at) VALUES ('guest', datetime('now'));

CREATE TABLE IF NOT EXISTS short_links (
  id TEXT PRIMARY KEY,
  owner_id TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  target_url TEXT NOT NULL,
  expires_at TEXT,
  revoked_at TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS assets (
  id TEXT PRIMARY KEY,
  owner_id TEXT NOT NULL,
  kind TEXT NOT NULL,
  name TEXT NOT NULL,
  content_type TEXT NOT NULL,
  size INTEGER NOT NULL,
  object_key TEXT NOT NULL,
  short_slug TEXT,
  expires_at TEXT,
  deleted_at TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS clipboard_items (
  id TEXT PRIMARY KEY,
  owner_id TEXT NOT NULL,
  content TEXT NOT NULL,
  password_hash TEXT NOT NULL DEFAULT '',
  short_slug TEXT,
  max_visits INTEGER NOT NULL DEFAULT 0,
  visit_count INTEGER NOT NULL DEFAULT 0,
  expires_at TEXT,
  deleted_at TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_short_links_slug ON short_links(slug);
CREATE INDEX IF NOT EXISTS idx_short_links_expires_at ON short_links(expires_at);
CREATE INDEX IF NOT EXISTS idx_assets_owner ON assets(owner_id);
CREATE INDEX IF NOT EXISTS idx_assets_kind ON assets(kind);
CREATE INDEX IF NOT EXISTS idx_assets_expires_at ON assets(expires_at);
CREATE INDEX IF NOT EXISTS idx_clipboard_expires_at ON clipboard_items(expires_at);
`
