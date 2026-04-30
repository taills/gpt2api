package sqltime
// Package sqltime provides a NullTime type that works correctly with the
// modernc.org/sqlite driver, which returns DATETIME column values as Go
// strings rather than time.Time. The standard sql.NullTime.Scan only accepts
// time.Time or nil, so any non-NULL DATETIME column would produce:
//
//	sql: Scan error … unsupported Scan, storing driver.Value type string into type *time.Time
//
// NullTime is a drop-in replacement for sql.NullTime: it implements
// sql.Scanner (handles nil, string, and time.Time inputs),
// driver.Valuer (writes back correctly), and json.Marshaler/Unmarshaler
// (serialises as a plain RFC3339 string or JSON null).
package sqltime

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// sqliteFmts is the set of datetime string formats the modernc.org/sqlite
// driver can return for a DATETIME/TIMESTAMP column.
var sqliteFmts = []string{
	time.RFC3339Nano,           // "2006-01-02T15:04:05.999999999Z07:00"
	time.RFC3339,               // "2006-01-02T15:04:05Z07:00"
	"2006-01-02T15:04:05",      // without timezone
	"2006-01-02 15:04:05.999999999-07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05",      // most common SQLite text datetime
	"2006-01-02",               // date-only
}

// NullTime represents a nullable time value. It is a drop-in replacement
// for database/sql.NullTime that can also scan string values returned by
// the modernc.org/sqlite driver.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true when Time is not NULL
}

// Scan implements sql.Scanner. It accepts nil (→ not valid), time.Time
// (→ direct), and string (→ parsed using known SQLite datetime formats).
func (n *NullTime) Scan(value any) error {
	if value == nil {
		n.Time, n.Valid = time.Time{}, false
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		n.Time, n.Valid = v, true
		return nil
	case string:
		if v == "" {
			n.Time, n.Valid = time.Time{}, false
			return nil
		}
		for _, layout := range sqliteFmts {
			t, err := time.Parse(layout, v)
			if err == nil {
				n.Time, n.Valid = t, true
				return nil
			}
		}
		return fmt.Errorf("sqltime: cannot parse %q as datetime", v)
	case []byte:
		return n.Scan(string(v))
	default:
		return fmt.Errorf("sqltime: unsupported type %T", value)
	}
}

// Value implements driver.Valuer so NullTime can be written back to the DB.
func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

// MarshalJSON serialises NullTime as a RFC3339Nano string or JSON null.
func (n NullTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Time)
}

// UnmarshalJSON parses a RFC3339 string or JSON null.
func (n *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &n.Time); err != nil {
		return err
	}
	n.Valid = true
	return nil
}

// From returns a valid NullTime wrapping t. If t is the zero value,
// the returned NullTime is still valid; callers should gate on !t.IsZero()
// when that matters.
func From(t time.Time) NullTime {
	return NullTime{Time: t, Valid: true}
}
