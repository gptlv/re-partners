package migrations

import "embed"

// Files exposes the embedded SQL migrations.
//
//go:embed *.sql
var Files embed.FS

// Dir is the top-level directory used when running goose.
const Dir = "."
