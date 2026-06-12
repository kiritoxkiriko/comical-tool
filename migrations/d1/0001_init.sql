CREATE TABLE IF NOT EXISTS short_links (
  id TEXT PRIMARY KEY,
  slug TEXT NOT NULL UNIQUE,
  target_url TEXT NOT NULL,
  expires_at INTEGER,
  revoked_at INTEGER,
  created_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS assets (
  id TEXT PRIMARY KEY,
  kind TEXT NOT NULL,
  name TEXT NOT NULL,
  content_type TEXT NOT NULL,
  size INTEGER NOT NULL,
  object_key TEXT NOT NULL,
  short_slug TEXT,
  password_hash TEXT NOT NULL DEFAULT '',
  max_visits INTEGER NOT NULL DEFAULT 0,
  visit_count INTEGER NOT NULL DEFAULT 0,
  expires_at INTEGER,
  deleted_at INTEGER,
  created_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_short_links_slug ON short_links(slug);
CREATE INDEX IF NOT EXISTS idx_assets_kind ON assets(kind);
CREATE INDEX IF NOT EXISTS idx_assets_expires_at ON assets(expires_at);
