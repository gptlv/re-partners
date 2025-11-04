package app

import "errors"

var (
	ErrSizeExists   = errors.New("pack size already exists")
	ErrSizeNotFound = errors.New("pack size not found")
	ErrLastSize     = errors.New("cannot delete the last pack size")
)
