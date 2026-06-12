CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  created_at INTEGER NOT NULL
);

INSERT OR IGNORE INTO users (id, created_at) VALUES ('guest', unixepoch());

ALTER TABLE short_links ADD COLUMN owner_id TEXT NOT NULL DEFAULT 'guest';

ALTER TABLE assets ADD COLUMN owner_id TEXT NOT NULL DEFAULT 'guest';

CREATE TABLE IF NOT EXISTS clipboard_items (
  id TEXT PRIMARY KEY,
  owner_id TEXT NOT NULL DEFAULT 'guest',
  content TEXT NOT NULL,
  password_hash TEXT NOT NULL DEFAULT '',
  short_slug TEXT,
  max_visits INTEGER NOT NULL DEFAULT 0,
  visit_count INTEGER NOT NULL DEFAULT 0,
  expires_at INTEGER,
  deleted_at INTEGER,
  created_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_assets_owner ON assets(owner_id);
CREATE INDEX IF NOT EXISTS idx_clipboard_expires_at ON clipboard_items(expires_at);
