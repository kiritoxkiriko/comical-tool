package repository

import (
	"context"
	"database/sql"

	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

type assetRow struct {
	ID          string         `db:"id"`
	Kind        string         `db:"kind"`
	Name        string         `db:"name"`
	ContentType string         `db:"content_type"`
	Size        int64          `db:"size"`
	ObjectKey   string         `db:"object_key"`
	ShortSlug   sql.NullString `db:"short_slug"`
	ExpiresAt   sql.NullString `db:"expires_at"`
	DeletedAt   sql.NullString `db:"deleted_at"`
	CreatedAt   string         `db:"created_at"`
}

func (s *SQLite) CreateAsset(ctx context.Context, asset domain.Asset) error {
	now := nowString()
	_, err := s.db.ExecContext(ctx, `
INSERT INTO assets
(id, owner_id, kind, name, content_type, size, object_key, short_slug, expires_at, deleted_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		asset.ID, domain.GuestUserID, asset.Kind, asset.Name, asset.ContentType,
		asset.Size, asset.ObjectKey, nullString(asset.ShortSlug),
		nullableTime(asset.ExpiresAt), nullableTime(asset.DeletedAt), now, now)
	return err
}

func (s *SQLite) ListAssets(ctx context.Context, kind domain.ResourceType) ([]domain.Asset, error) {
	rows := []assetRow{}
	err := s.db.SelectContext(ctx, &rows, `
SELECT id, kind, name, content_type, size, object_key, short_slug, expires_at, deleted_at, created_at
FROM assets WHERE kind = ? AND deleted_at IS NULL ORDER BY created_at DESC`, kind)
	if err != nil {
		return nil, err
	}
	assets := make([]domain.Asset, 0, len(rows))
	for _, row := range rows {
		assets = append(assets, row.toDomain())
	}
	return assets, nil
}

func (s *SQLite) FindAsset(ctx context.Context, id string) (domain.Asset, error) {
	var row assetRow
	err := s.db.GetContext(ctx, &row, `
SELECT id, kind, name, content_type, size, object_key, short_slug, expires_at, deleted_at, created_at
FROM assets WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return domain.Asset{}, ErrNotFound
	}
	if err != nil {
		return domain.Asset{}, err
	}
	return row.toDomain(), nil
}

func (s *SQLite) DeleteAsset(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `
UPDATE assets SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`,
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

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func (r assetRow) toDomain() domain.Asset {
	return domain.Asset{
		ID:          r.ID,
		Kind:        domain.ResourceType(r.Kind),
		Name:        r.Name,
		ContentType: r.ContentType,
		Size:        r.Size,
		ObjectKey:   r.ObjectKey,
		ShortSlug:   r.ShortSlug.String,
		ExpiresAt:   parseNullableTime(r.ExpiresAt),
		DeletedAt:   parseNullableTime(r.DeletedAt),
		CreatedAt:   parseTime(r.CreatedAt),
	}
}
