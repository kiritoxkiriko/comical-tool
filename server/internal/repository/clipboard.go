package repository

import (
	"context"
	"database/sql"

	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

type clipboardRow struct {
	ID           string         `db:"id"`
	Content      string         `db:"content"`
	PasswordHash string         `db:"password_hash"`
	ShortSlug    sql.NullString `db:"short_slug"`
	MaxVisits    int            `db:"max_visits"`
	VisitCount   int            `db:"visit_count"`
	ExpiresAt    sql.NullString `db:"expires_at"`
	DeletedAt    sql.NullString `db:"deleted_at"`
	CreatedAt    string         `db:"created_at"`
}

func (s *SQLite) CreateClipboard(ctx context.Context, item domain.ClipboardItem) error {
	now := nowString()
	_, err := s.db.ExecContext(ctx, `
INSERT INTO clipboard_items
(id, owner_id, content, password_hash, short_slug, max_visits, visit_count, expires_at, deleted_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, domain.GuestUserID, item.Content, item.PasswordHash,
		nullString(item.ShortSlug), item.MaxVisits, item.VisitCount,
		nullableTime(item.ExpiresAt), nullableTime(item.DeletedAt), now, now)
	return err
}

func (s *SQLite) FindClipboard(ctx context.Context, id string) (domain.ClipboardItem, error) {
	var row clipboardRow
	err := s.db.GetContext(ctx, &row, `
SELECT id, content, password_hash, short_slug, max_visits, visit_count, expires_at, deleted_at, created_at
FROM clipboard_items WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return domain.ClipboardItem{}, ErrNotFound
	}
	if err != nil {
		return domain.ClipboardItem{}, err
	}
	return row.toDomain(), nil
}

func (s *SQLite) DeleteClipboard(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `
UPDATE clipboard_items SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`,
		nowString(), nowString(), id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLite) IncrementClipboardVisit(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `
UPDATE clipboard_items SET visit_count = visit_count + 1, updated_at = ? WHERE id = ?`,
		nowString(), id)
	return err
}

func (r clipboardRow) toDomain() domain.ClipboardItem {
	return domain.ClipboardItem{
		ID:           r.ID,
		Content:      r.Content,
		PasswordHash: r.PasswordHash,
		ShortSlug:    r.ShortSlug.String,
		MaxVisits:    r.MaxVisits,
		VisitCount:   r.VisitCount,
		ExpiresAt:    parseNullableTime(r.ExpiresAt),
		DeletedAt:    parseNullableTime(r.DeletedAt),
		CreatedAt:    parseTime(r.CreatedAt),
	}
}
