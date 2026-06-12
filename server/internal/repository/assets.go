package repository

import (
	"context"
	"database/sql"

	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

type assetRow struct {
	ID           string         `db:"id"`
	Kind         string         `db:"kind"`
	Name         string         `db:"name"`
	ContentType  string         `db:"content_type"`
	Size         int64          `db:"size"`
	ObjectKey    string         `db:"object_key"`
	ShortSlug    sql.NullString `db:"short_slug"`
	PasswordHash string         `db:"password_hash"`
	MaxVisits    int            `db:"max_visits"`
	VisitCount   int            `db:"visit_count"`
	ExpiresAt    dbTime         `db:"expires_at"`
	DeletedAt    dbTime         `db:"deleted_at"`
	CreatedAt    dbTime         `db:"created_at"`
}

func (s *Store) CreateAsset(ctx context.Context, asset domain.Asset) error {
	now := s.nowArg()
	_, err := s.exec(ctx, `
INSERT INTO assets
(id, owner_id, kind, name, content_type, size, object_key, short_slug, password_hash, max_visits, visit_count, expires_at, deleted_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		asset.ID, domain.GuestUserID, asset.Kind, asset.Name, asset.ContentType,
		asset.Size, asset.ObjectKey, nullString(asset.ShortSlug), asset.PasswordHash, asset.MaxVisits, asset.VisitCount,
		s.timeArg(nullableTime(asset.ExpiresAt)), s.timeArg(nullableTime(asset.DeletedAt)), now, now)
	return err
}

func (s *Store) ListAssets(ctx context.Context, kind domain.ResourceType) ([]domain.Asset, error) {
	rows := []assetRow{}
	err := s.selectRows(ctx, &rows, `
SELECT id, kind, name, content_type, size, object_key, short_slug, password_hash, max_visits, visit_count, expires_at, deleted_at, created_at
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

func (s *Store) FindAsset(ctx context.Context, id string) (domain.Asset, error) {
	var row assetRow
	err := s.get(ctx, &row, `
SELECT id, kind, name, content_type, size, object_key, short_slug, password_hash, max_visits, visit_count, expires_at, deleted_at, created_at
FROM assets WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return domain.Asset{}, ErrNotFound
	}
	if err != nil {
		return domain.Asset{}, err
	}
	return row.toDomain(), nil
}

func (s *Store) DeleteAsset(ctx context.Context, id string) error {
	now := s.nowArg()
	res, err := s.exec(ctx, `
UPDATE assets SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`,
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

func (s *Store) IncrementAssetVisit(ctx context.Context, id string) error {
	_, err := s.exec(ctx, "UPDATE assets SET visit_count = visit_count + 1, updated_at = ? WHERE id = ?", s.nowArg(), id)
	return err
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func (r assetRow) toDomain() domain.Asset {
	return domain.Asset{
		ID:           r.ID,
		Kind:         domain.ResourceType(r.Kind),
		Name:         r.Name,
		ContentType:  r.ContentType,
		Size:         r.Size,
		ObjectKey:    r.ObjectKey,
		ShortSlug:    r.ShortSlug.String,
		PasswordHash: r.PasswordHash,
		MaxVisits:    r.MaxVisits,
		VisitCount:   r.VisitCount,
		ExpiresAt:    parseNullableTime(r.ExpiresAt),
		DeletedAt:    parseNullableTime(r.DeletedAt),
		CreatedAt:    parseTime(r.CreatedAt),
	}
}
