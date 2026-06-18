// Package migrations embeds the SQL migration files so they ship inside the
// compiled binary and can be applied automatically on startup.
package migrations

import "embed"

// FS holds all *.sql migration files in this directory.
//
//go:embed *.sql
var FS embed.FS
