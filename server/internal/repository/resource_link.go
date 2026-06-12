package repository

import (
	"context"

	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

type resourceLinkRow struct {
	ID           string `db:"id"`
	ShortLinkID  string `db:"short_link_id"`
	ResourceType string `db:"resource_type"`
	ResourceID   string `db:"resource_id"`
	CreatedAt    dbTime `db:"created_at"`
}

func (s *Store) CreateResourceLink(ctx context.Context, link domain.ResourceLink) error {
	_, err := s.exec(ctx, `
INSERT INTO resource_links
(id, short_link_id, resource_type, resource_id, created_at)
VALUES (?, ?, ?, ?, ?)`,
		link.ID, link.ShortLinkID, link.ResourceType, link.ResourceID, s.nowArg())
	return err
}

func (s *Store) ListResourceLinks(ctx context.Context, resourceType domain.ResourceType, resourceID string) ([]domain.ResourceLink, error) {
	rows := []resourceLinkRow{}
	err := s.selectRows(ctx, &rows, `
SELECT id, short_link_id, resource_type, resource_id, created_at
FROM resource_links WHERE resource_type = ? AND resource_id = ? ORDER BY created_at DESC`,
		resourceType, resourceID)
	if err != nil {
		return nil, err
	}
	links := make([]domain.ResourceLink, 0, len(rows))
	for _, row := range rows {
		links = append(links, row.toDomain())
	}
	return links, nil
}

func (r resourceLinkRow) toDomain() domain.ResourceLink {
	return domain.ResourceLink{
		ID:           r.ID,
		ShortLinkID:  r.ShortLinkID,
		ResourceType: domain.ResourceType(r.ResourceType),
		ResourceID:   r.ResourceID,
		CreatedAt:    parseTime(r.CreatedAt),
	}
}
