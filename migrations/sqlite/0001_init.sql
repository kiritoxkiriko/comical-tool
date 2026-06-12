CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  created_at TEXT NOT NULL
);

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
  password_hash TEXT NOT NULL DEFAULT '',
  max_visits INTEGER NOT NULL DEFAULT 0,
  visit_count INTEGER NOT NULL DEFAULT 0,
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
