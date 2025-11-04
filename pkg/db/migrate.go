package db

import (
	"database/sql"
	"fmt"

	"github.com/gptlv/re-partners/packs/migrations"
	"github.com/pressly/goose/v3"
)

// RunMigrations applies all embedded goose migrations against the provided database.
func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrations.Files)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("configure goose dialect: %w", err)
	}

	if err := goose.Up(db, migrations.Dir); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
