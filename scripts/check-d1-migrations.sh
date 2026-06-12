#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d)"
DB_PATH="${TMP_DIR}/d1.sqlite"

cleanup() {
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

for migration in "${ROOT_DIR}"/migrations/d1/*.sql; do
  sqlite3 "${DB_PATH}" < "${migration}"
done

require_table() {
  local table="$1"
  local count
  count="$(sqlite3 "${DB_PATH}" "SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = '${table}';")"
  if [[ "${count}" != "1" ]]; then
    echo "missing D1 table: ${table}" >&2
    exit 1
  fi
}

require_index() {
  local index="$1"
  local count
  count="$(sqlite3 "${DB_PATH}" "SELECT count(*) FROM sqlite_master WHERE type = 'index' AND name = '${index}';")"
  if [[ "${count}" != "1" ]]; then
    echo "missing D1 index: ${index}" >&2
    exit 1
  fi
}

require_column() {
  local table="$1"
  local column="$2"
  local count
  count="$(sqlite3 "${DB_PATH}" "SELECT count(*) FROM pragma_table_info('${table}') WHERE name = '${column}';")"
  if [[ "${count}" != "1" ]]; then
    echo "missing D1 column: ${table}.${column}" >&2
    exit 1
  fi
}

require_table users
require_table short_links
require_table assets
require_table clipboard_items
require_table access_events
require_table resource_links

require_column users id
require_column users created_at

require_column short_links owner_id
require_column assets owner_id

require_column clipboard_items id
require_column clipboard_items owner_id
require_column clipboard_items content
require_column clipboard_items password_hash
require_column clipboard_items short_slug
require_column clipboard_items max_visits
require_column clipboard_items visit_count
require_column clipboard_items expires_at
require_column clipboard_items deleted_at
require_column clipboard_items created_at

require_index idx_assets_owner
require_index idx_clipboard_expires_at

guest_count="$(sqlite3 "${DB_PATH}" "SELECT count(*) FROM users WHERE id = 'guest';")"
if [[ "${guest_count}" != "1" ]]; then
  echo "missing D1 guest user seed" >&2
  exit 1
fi

echo "D1 migrations OK"
