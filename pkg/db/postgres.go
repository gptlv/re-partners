package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open returns a PostgreSQL-backed sql.DB using the pgx stdlib driver.
func Open(connString string) (*sql.DB, error) {
	return sql.Open("pgx", connString)
}
