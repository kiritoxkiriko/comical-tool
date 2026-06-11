package repository

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const timeLayout = time.RFC3339Nano

type dbTime struct {
	Time  time.Time
	Valid bool
}

func nullableTime(t *time.Time) dbTime {
	if t == nil {
		return dbTime{}
	}
	return newDBTime(*t)
}

func newDBTime(value time.Time) dbTime {
	return dbTime{Time: value.UTC(), Valid: true}
}

func (t *dbTime) Scan(value any) error {
	if value == nil {
		*t = dbTime{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = newDBTime(v)
		return nil
	case string:
		return t.scanString(v)
	case []byte:
		return t.scanString(string(v))
	default:
		return fmt.Errorf("unsupported time value %T", value)
	}
}

func (t dbTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.UTC(), nil
}

func (t *dbTime) scanString(value string) error {
	if value == "" {
		*t = dbTime{}
		return nil
	}
	for _, layout := range []string{timeLayout, time.RFC3339, "2006-01-02 15:04:05.999999999-07:00", "2006-01-02 15:04:05.999999", "2006-01-02 15:04:05"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			*t = newDBTime(parsed)
			return nil
		}
	}
	return fmt.Errorf("invalid time value %q", value)
}

func (t dbTime) ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	utc := t.Time.UTC()
	return &utc
}

func parseTime(value dbTime) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time
}

func parseNullableTime(value dbTime) *time.Time {
	return value.ptr()
}

func now() time.Time {
	return time.Now().UTC()
}
