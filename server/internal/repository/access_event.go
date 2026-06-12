package repository

import (
	"context"

	"github.com/kiritoxkiriko/comical-tool/server/pkg/domain"
)

type accessEventRow struct {
	ID           string `db:"id"`
	ResourceType string `db:"resource_type"`
	ResourceID   string `db:"resource_id"`
	Action       string `db:"action"`
	CreatedAt    dbTime `db:"created_at"`
}

func (s *Store) RecordAccessEvent(ctx context.Context, event domain.AccessEvent) error {
	_, err := s.exec(ctx, `
INSERT INTO access_events
(id, resource_type, resource_id, action, created_at)
VALUES (?, ?, ?, ?, ?)`,
		event.ID, event.ResourceType, event.ResourceID, event.Action, s.nowArg())
	return err
}

func (s *Store) ListAccessEvents(ctx context.Context, resourceType domain.ResourceType, resourceID string) ([]domain.AccessEvent, error) {
	rows := []accessEventRow{}
	err := s.selectRows(ctx, &rows, `
SELECT id, resource_type, resource_id, action, created_at
FROM access_events WHERE resource_type = ? AND resource_id = ? ORDER BY created_at DESC`,
		resourceType, resourceID)
	if err != nil {
		return nil, err
	}
	events := make([]domain.AccessEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, row.toDomain())
	}
	return events, nil
}

func (r accessEventRow) toDomain() domain.AccessEvent {
	return domain.AccessEvent{
		ID:           r.ID,
		ResourceType: domain.ResourceType(r.ResourceType),
		ResourceID:   r.ResourceID,
		Action:       r.Action,
		CreatedAt:    parseTime(r.CreatedAt),
	}
}
