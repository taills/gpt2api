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

// splitStatements splits a SQL file on ";" delimiters, respecting single-quoted
// string literals, and trims whitespace / comment-only entries.
func splitStatements(sql string) []string {
	var stmts []string
	var buf strings.Builder
	inSingleQuote := false

	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		switch {
		case ch == '\'' && !inSingleQuote:
			inSingleQuote = true
			buf.WriteByte(ch)
		case ch == '\'' && inSingleQuote:
			buf.WriteByte(ch)
			// handle escaped single-quote ''
			if i+1 < len(sql) && sql[i+1] == '\'' {
				i++
				buf.WriteByte(sql[i])
			} else {
				inSingleQuote = false
			}
		case ch == ';' && !inSingleQuote:
			stmt := strings.TrimSpace(buf.String())
			buf.Reset()
			if stmt == "" {
				continue
			}
			// Skip entries that consist only of comment lines
			allComment := true
			for _, line := range strings.Split(stmt, "\n") {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "--") {
					allComment = false
					break
				}
			}
			if !allComment {
				stmts = append(stmts, stmt)
			}
		default:
			buf.WriteByte(ch)
		}
	}
	// Handle trailing statement without a trailing semicolon
	if stmt := strings.TrimSpace(buf.String()); stmt != "" {
		allComment := true
		for _, line := range strings.Split(stmt, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "--") {
				allComment = false
				break
			}
		}
		if !allComment {
			stmts = append(stmts, stmt)
		}
	}
	return stmts
}
