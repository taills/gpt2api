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
	"strings"
	"time"
)

// sqliteFmts is the set of datetime string formats the modernc.org/sqlite
// driver can return for a DATETIME/TIMESTAMP column.
//
// The format "2006-01-02 15:04:05 -0700 MST" (and its nanosecond variant)
// handles values stored via Go's time.Time.String(), which the modernc.org/sqlite
// driver produces when a time.Time driver.Value is written without explicit
// formatting (e.g. "2026-05-03 14:54:33 +0800 CST").
var sqliteFmts = []string{
	time.RFC3339Nano,      // "2006-01-02T15:04:05.999999999Z07:00"
	time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
	"2006-01-02T15:04:05", // without timezone
	"2006-01-02 15:04:05.999999999-07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05",                     // most common SQLite text datetime
	"2006-01-02 15:04:05.999999999 -0700 MST", // Go time.String() with nanoseconds + alpha zone
	"2006-01-02 15:04:05 -0700 MST",           // Go time.String() without nanoseconds + alpha zone
	"2006-01-02 15:04:05.999999999 -0700",     // after stripNumericZoneName, with nanoseconds
	"2006-01-02 15:04:05 -0700",               // after stripNumericZoneName: "2026-05-10 18:54:12 +0800"
	"2006-01-02",                              // date-only
}

// NullTime represents a nullable time value. It is a drop-in replacement
// for database/sql.NullTime that can also scan string values returned by
// the modernc.org/sqlite driver.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true when Time is not NULL
}

// stripNumericZoneName strips the trailing numeric zone-name token that Go's
// time.String() appends when the timezone was created via time.FixedZone with a
// numeric name (e.g. "+0800").  In that case time.String() produces strings like
//
//	"2026-05-10 18:54:12 +0800 +0800"
//
// where the last token is the zone name ("+0800") rather than an alphabetic
// abbreviation ("CST").  time.Parse's "MST" verb only accepts alphabetic zone
// abbreviations, so we strip the numeric name before parsing.
//
// The function is a no-op when the trailing token is alphabetic or when the
// string does not match the expected shape.
func stripNumericZoneName(s string) string {
	// Find the last space-separated token.
	lastSpace := strings.LastIndex(s, " ")
	if lastSpace < 0 {
		return s
	}
	tail := s[lastSpace+1:]
	// A numeric zone-name looks like "+0800" or "-0700": sign + 4 digits.
	if len(tail) != 5 {
		return s
	}
	if tail[0] != '+' && tail[0] != '-' {
		return s
	}
	for _, c := range tail[1:] {
		if c < '0' || c > '9' {
			return s // alphabetic abbreviation – leave it
		}
	}
	// Strip the numeric zone-name; the offset before it is sufficient.
	return strings.TrimSpace(s[:lastSpace])
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
		// Normalise strings produced by time.String() with a FixedZone whose name
		// is a numeric offset (e.g. "2026-05-10 18:54:12 +0800 +0800").
		normalized := stripNumericZoneName(v)
		for _, layout := range sqliteFmts {
			t, err := time.Parse(layout, normalized)
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
// Returns a UTC RFC3339Nano string rather than time.Time to prevent the
// modernc.org/sqlite driver from using Go's time.String() format
// (e.g. "2026-05-03 14:54:33 +0800 CST"), which cannot be read back.
func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time.UTC().Format(time.RFC3339Nano), nil
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
