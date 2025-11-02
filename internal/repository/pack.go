package repository

import (
	"context"
	"database/sql"
)

type PackRepository struct {
	db *sql.DB
}

func NewPackRepository(db *sql.DB) *PackRepository {
	return &PackRepository{db: db}
}

// Sizes returns all configured pack sizes sorted ascending.
func (r *PackRepository) Sizes(ctx context.Context) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT size
		FROM pack_sizes
		ORDER BY size ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sizes []int64
	for rows.Next() {
		var size int64
		if err := rows.Scan(&size); err != nil {
			return nil, err
		}
		sizes = append(sizes, size)
	}

	return sizes, rows.Err()
}
