package policy

import "time"

func ExpiryFromDuration(ttl time.Duration) *time.Time {
	if ttl <= 0 {
		return nil
	}
	expiresAt := time.Now().UTC().Add(ttl)
	return &expiresAt
}

func IsExpired(expiresAt *time.Time) bool {
	return expiresAt != nil && time.Now().UTC().After(*expiresAt)
}
