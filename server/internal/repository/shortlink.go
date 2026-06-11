package repository

import (
	"context"
	"database/sql"

	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

type shortRow struct {
	ID        string         `db:"id"`
	Slug      string         `db:"slug"`
	TargetURL string         `db:"target_url"`
	ExpiresAt sql.NullString `db:"expires_at"`
	RevokedAt sql.NullString `db:"revoked_at"`
	CreatedAt string         `db:"created_at"`
}

func (s *SQLite) CreateShortLink(ctx context.Context, link domain.ShortLink) error {
	now := nowString()
	_, err := s.db.ExecContext(ctx, `
INSERT INTO short_links
(id, owner_id, slug, target_url, expires_at, revoked_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		link.ID, domain.GuestUserID, link.Slug, link.TargetURL,
		nullableTime(link.ExpiresAt), nullableTime(link.RevokedAt), now, now)
	return err
}

func (s *SQLite) FindShortLink(ctx context.Context, slug string) (domain.ShortLink, error) {
	var row shortRow
	err := s.db.GetContext(ctx, &row, `
SELECT id, slug, target_url, expires_at, revoked_at, created_at
FROM short_links WHERE slug = ?`, slug)
	if err == sql.ErrNoRows {
		return domain.ShortLink{}, ErrNotFound
	}
	if err != nil {
		return domain.ShortLink{}, err
	}
	return row.toDomain(), nil
}

func (s *SQLite) RevokeShortLink(ctx context.Context, slug string) error {
	res, err := s.db.ExecContext(ctx, `
UPDATE short_links SET revoked_at = ?, updated_at = ? WHERE slug = ?`,
		nowString(), nowString(), slug)
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

func (r shortRow) toDomain() domain.ShortLink {
	return domain.ShortLink{
		ID:        r.ID,
		Slug:      r.Slug,
		TargetURL: r.TargetURL,
		ExpiresAt: parseNullableTime(r.ExpiresAt),
		RevokedAt: parseNullableTime(r.RevokedAt),
		CreatedAt: parseTime(r.CreatedAt),
	}
}
