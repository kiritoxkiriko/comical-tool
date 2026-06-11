CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(191) PRIMARY KEY,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
);

INSERT IGNORE INTO users (id) VALUES ('guest');

CREATE TABLE IF NOT EXISTS short_links (
  id VARCHAR(191) PRIMARY KEY,
  owner_id VARCHAR(191) NOT NULL DEFAULT 'guest',
  slug VARCHAR(191) NOT NULL UNIQUE,
  target_url TEXT NOT NULL,
  expires_at DATETIME(3),
  revoked_at DATETIME(3),
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
);

CREATE TABLE IF NOT EXISTS assets (
  id VARCHAR(191) PRIMARY KEY,
  owner_id VARCHAR(191) NOT NULL DEFAULT 'guest',
  kind VARCHAR(32) NOT NULL,
  name TEXT NOT NULL,
  content_type VARCHAR(255) NOT NULL,
  size BIGINT NOT NULL,
  object_key TEXT NOT NULL,
  short_slug VARCHAR(191),
  expires_at DATETIME(3),
  deleted_at DATETIME(3),
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
);

CREATE TABLE IF NOT EXISTS clipboard_items (
  id VARCHAR(191) PRIMARY KEY,
  owner_id VARCHAR(191) NOT NULL DEFAULT 'guest',
  content LONGTEXT NOT NULL,
  password_hash TEXT NOT NULL,
  short_slug VARCHAR(191),
  max_visits INT NOT NULL DEFAULT 0,
  visit_count INT NOT NULL DEFAULT 0,
  expires_at DATETIME(3),
  deleted_at DATETIME(3),
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
);

CREATE INDEX idx_short_links_slug ON short_links(slug);
CREATE INDEX idx_short_links_expires_at ON short_links(expires_at);
CREATE INDEX idx_assets_owner ON assets(owner_id);
CREATE INDEX idx_assets_kind ON assets(kind);
CREATE INDEX idx_assets_expires_at ON assets(expires_at);
CREATE INDEX idx_clipboard_expires_at ON clipboard_items(expires_at);
