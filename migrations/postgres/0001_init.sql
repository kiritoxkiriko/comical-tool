CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO users (id) VALUES ('guest')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS short_links (
  id TEXT PRIMARY KEY,
  owner_id TEXT NOT NULL DEFAULT 'guest',
  slug TEXT NOT NULL UNIQUE,
  target_url TEXT NOT NULL,
  expires_at TIMESTAMPTZ,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS assets (
  id TEXT PRIMARY KEY,
  owner_id TEXT NOT NULL DEFAULT 'guest',
  kind TEXT NOT NULL,
  name TEXT NOT NULL,
  content_type TEXT NOT NULL,
  size BIGINT NOT NULL,
  object_key TEXT NOT NULL,
  short_slug TEXT,
  password_hash TEXT NOT NULL DEFAULT '',
  max_visits INTEGER NOT NULL DEFAULT 0,
  visit_count INTEGER NOT NULL DEFAULT 0,
  expires_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS clipboard_items (
  id TEXT PRIMARY KEY,
  owner_id TEXT NOT NULL DEFAULT 'guest',
  content TEXT NOT NULL,
  password_hash TEXT NOT NULL DEFAULT '',
  short_slug TEXT,
  max_visits INTEGER NOT NULL DEFAULT 0,
  visit_count INTEGER NOT NULL DEFAULT 0,
  expires_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_short_links_slug ON short_links(slug);
CREATE INDEX IF NOT EXISTS idx_short_links_expires_at ON short_links(expires_at);
CREATE INDEX IF NOT EXISTS idx_assets_owner ON assets(owner_id);
CREATE INDEX IF NOT EXISTS idx_assets_kind ON assets(kind);
CREATE INDEX IF NOT EXISTS idx_assets_expires_at ON assets(expires_at);
CREATE INDEX IF NOT EXISTS idx_clipboard_expires_at ON clipboard_items(expires_at);
