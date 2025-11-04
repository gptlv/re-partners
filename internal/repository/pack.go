package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sqlite3 "modernc.org/sqlite/lib"
)

var (
	ErrDuplicateSize = errors.New("duplicate pack size")
	ErrSizeNotFound  = errors.New("pack size not found")
)

type PackSize struct {
	ID   int64
	Size int64
}

type PackRepository struct {
	db *sql.DB
}

func NewPackRepository(db *sql.DB) *PackRepository {
	return &PackRepository{db: db}
}

// Sizes returns all configured pack sizes sorted ascending.
func (r *PackRepository) Sizes(ctx context.Context) ([]PackSize, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, size
		FROM pack_sizes
		ORDER BY size ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sizes []PackSize
	for rows.Next() {
		var ps PackSize
		if err := rows.Scan(&ps.ID, &ps.Size); err != nil {
			return nil, err
		}
		sizes = append(sizes, ps)
	}

	return sizes, rows.Err()
}

// AddSize inserts a new pack size record.
func (r *PackRepository) AddSize(ctx context.Context, size int64) (*PackSize, error) {
	var ps PackSize
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO pack_sizes (size)
		VALUES (?)
		RETURNING id, size
	`, size).Scan(&ps.ID, &ps.Size)
	if err != nil {
		if isConstraintError(err) {
			return nil, fmt.Errorf("%w: %v", ErrDuplicateSize, err)
		}
		return nil, err
	}

	return &ps, nil
}

// DeleteSize removes a pack size by ID.
func (r *PackRepository) DeleteSize(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM pack_sizes
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrSizeNotFound
	}

	return nil
}

// CountSizes returns the number of configured pack sizes.
func (r *PackRepository) CountSizes(ctx context.Context) (int, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM pack_sizes
	`).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// EnsureSizeExists returns ErrSizeNotFound if the pack size is missing.
func (r *PackRepository) EnsureSizeExists(ctx context.Context, id int64) error {
	var exists bool
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(1) > 0
		FROM pack_sizes
		WHERE id = ?
	`, id).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return ErrSizeNotFound
	}
	return nil
}

func isConstraintError(err error) bool {
	var sqliteErr interface {
		error
		Code() int
	}
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT
	}
	return false
}
