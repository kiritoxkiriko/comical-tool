package repository

import (
	"database/sql"
	"time"
)

const timeLayout = time.RFC3339Nano

func nullableTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(timeLayout)
}

func parseTime(value string) time.Time {
	parsed, _ := time.Parse(timeLayout, value)
	return parsed
}

func parseNullableTime(value sql.NullString) *time.Time {
	if !value.Valid || value.String == "" {
		return nil
	}
	parsed := parseTime(value.String)
	return &parsed
}

func nowString() string {
	return time.Now().UTC().Format(timeLayout)
}
