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
	ExpiresAt    dbTime         `db:"expires_at"`
	DeletedAt    dbTime         `db:"deleted_at"`
	CreatedAt    dbTime         `db:"created_at"`
}

func (s *Store) CreateClipboard(ctx context.Context, item domain.ClipboardItem) error {
	now := s.nowArg()
	_, err := s.exec(ctx, `
INSERT INTO clipboard_items
(id, owner_id, content, password_hash, short_slug, max_visits, visit_count, expires_at, deleted_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, domain.GuestUserID, item.Content, item.PasswordHash,
		nullString(item.ShortSlug), item.MaxVisits, item.VisitCount,
		s.timeArg(nullableTime(item.ExpiresAt)), s.timeArg(nullableTime(item.DeletedAt)), now, now)
	return err
}

func (s *Store) FindClipboard(ctx context.Context, id string) (domain.ClipboardItem, error) {
	var row clipboardRow
	err := s.get(ctx, &row, `
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

func (s *Store) DeleteClipboard(ctx context.Context, id string) error {
	now := s.nowArg()
	res, err := s.exec(ctx, `
UPDATE clipboard_items SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`,
		now, now, id)
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

func (s *Store) IncrementClipboardVisit(ctx context.Context, id string) error {
	_, err := s.exec(ctx, `
UPDATE clipboard_items SET visit_count = visit_count + 1, updated_at = ? WHERE id = ?`,
		s.nowArg(), id)
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
