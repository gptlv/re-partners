package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Open returns a SQLite-backed sql.DB using the modernc pure Go driver.
func Open(path string) (*sql.DB, error) {
	return sql.Open("sqlite", path)
}
