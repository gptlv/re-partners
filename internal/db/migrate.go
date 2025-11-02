package db

import (
	"context"
	"database/sql"
)

// Migrate ensures pack_sizes exists and seeds default rows.
func Migrate(db *sql.DB) error {
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS pack_sizes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	size INTEGER NOT NULL UNIQUE
);
`)
	if err != nil {
		return err
	}

	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pack_sizes`).Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		_, err = db.ExecContext(ctx, `
INSERT INTO pack_sizes (size) VALUES (250),(500),(1000),(2000),(5000);
`)
		if err != nil {
			return err
		}
	}

	return nil
}
