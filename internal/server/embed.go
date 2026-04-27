package server

import "embed"

// webFS holds the compiled frontend assets (web/dist/) embedded at build time.
// The build step copies web/dist/* into internal/server/web/ before `go build`.
// At development time the directory may be empty; spa.go falls back to disk.
//
//go:embed web
var webFS embed.FS
