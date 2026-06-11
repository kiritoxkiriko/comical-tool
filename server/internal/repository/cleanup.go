package repository

import "context"

type CleanupResult struct {
	Assets     int64 `json:"assets"`
	Clipboard  int64 `json:"clipboard"`
	ShortLinks int64 `json:"short_links"`
}

func (s *SQLite) CleanupExpired(ctx context.Context) (CleanupResult, error) {
	now := nowString()
	assets, err := s.db.ExecContext(ctx, `
UPDATE assets SET deleted_at = ?, updated_at = ? WHERE deleted_at IS NULL AND expires_at IS NOT NULL AND expires_at < ?`,
		now, now, now)
	if err != nil {
		return CleanupResult{}, err
	}
	clip, err := s.db.ExecContext(ctx, `
UPDATE clipboard_items SET deleted_at = ?, updated_at = ? WHERE deleted_at IS NULL AND expires_at IS NOT NULL AND expires_at < ?`,
		now, now, now)
	if err != nil {
		return CleanupResult{}, err
	}
	short, err := s.db.ExecContext(ctx, `
UPDATE short_links SET revoked_at = ?, updated_at = ? WHERE revoked_at IS NULL AND expires_at IS NOT NULL AND expires_at < ?`,
		now, now, now)
	if err != nil {
		return CleanupResult{}, err
	}
	return cleanupResult(assets, clip, short)
}

func cleanupResult(results ...interface{ RowsAffected() (int64, error) }) (CleanupResult, error) {
	counts := make([]int64, len(results))
	for idx, result := range results {
		count, err := result.RowsAffected()
		if err != nil {
			return CleanupResult{}, err
		}
		counts[idx] = count
	}
	return CleanupResult{Assets: counts[0], Clipboard: counts[1], ShortLinks: counts[2]}, nil
}
