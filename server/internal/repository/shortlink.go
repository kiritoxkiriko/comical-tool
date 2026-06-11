package repository

import (
	"context"
	"database/sql"

	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

type shortRow struct {
	ID        string `db:"id"`
	Slug      string `db:"slug"`
	TargetURL string `db:"target_url"`
	ExpiresAt dbTime `db:"expires_at"`
	RevokedAt dbTime `db:"revoked_at"`
	CreatedAt dbTime `db:"created_at"`
}

func (s *Store) CreateShortLink(ctx context.Context, link domain.ShortLink) error {
	now := s.nowArg()
	_, err := s.exec(ctx, `
INSERT INTO short_links
(id, owner_id, slug, target_url, expires_at, revoked_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		link.ID, domain.GuestUserID, link.Slug, link.TargetURL,
		s.timeArg(nullableTime(link.ExpiresAt)), s.timeArg(nullableTime(link.RevokedAt)), now, now)
	return err
}

func (s *Store) FindShortLink(ctx context.Context, slug string) (domain.ShortLink, error) {
	var row shortRow
	err := s.get(ctx, &row, `
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

func (s *Store) RevokeShortLink(ctx context.Context, slug string) error {
	now := s.nowArg()
	res, err := s.exec(ctx, `
UPDATE short_links SET revoked_at = ?, updated_at = ? WHERE slug = ?`,
		now, now, slug)
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
