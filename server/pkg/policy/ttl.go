package policy

import (
	"strconv"
	"strings"
	"time"
)

func ParseTTLDuration(value string, fallback time.Duration) (time.Duration, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback, nil
	}
	if strings.HasSuffix(value, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(value, "d"))
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	return time.ParseDuration(value)
}

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
