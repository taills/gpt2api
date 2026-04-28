package db

import (
	_ "embed"
	"strings"

	"github.com/jmoiron/sqlx"
)

//go:embed schema.sql
var schemaSql string

// Migrate runs schema.sql idempotently against db.
// Each SQL statement is executed individually; blank lines and
// comment-only lines are skipped.
func Migrate(db *sqlx.DB) error {
	for _, stmt := range splitStatements(schemaSql) {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

// splitStatements splits a SQL file on ";" delimiters and trims whitespace /
// comment-only entries.
func splitStatements(sql string) []string {
	parts := strings.Split(sql, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// Skip entries that consist only of comment lines
		allComment := true
		for _, line := range strings.Split(p, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "--") {
				allComment = false
				break
			}
		}
		if !allComment {
			out = append(out, p)
		}
	}
	return out
}
