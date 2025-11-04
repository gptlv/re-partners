package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/gptlv/re-partners/packs/internal/repository"
	"github.com/gptlv/re-partners/packs/pkg/calculate"
)

type Service struct {
	repo *repository.PackRepository
}

func NewService(repo *repository.PackRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Sizes(ctx context.Context) ([]PackSize, error) {
	rows, err := s.repo.Sizes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sizes: %w", err)
	}

	sizes := make([]PackSize, len(rows))
	for i, row := range rows {
		sizes[i] = PackSize{
			ID:   row.ID,
			Size: row.Size,
		}
	}

	return sizes, nil
}

func (s *Service) CalculatePackages(ctx context.Context, orderedItems int64) ([]Pack, error) {
	sizes, err := s.Sizes(ctx)
	if err != nil {
		return nil, err
	}

	if len(sizes) == 0 {
		return nil, fmt.Errorf("failed to calculate packages: %w", calculate.ErrEmptySizes)
	}

	sizeValues := make([]int64, len(sizes))
	for i, ps := range sizes {
		sizeValues[i] = ps.Size
	}

	result, err := calculate.CalculatePackages(orderedItems, sizeValues)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate packages: %w", err)
	}

	packs := make([]Pack, 0, len(result))
	for _, ps := range sizes {
		count, ok := result[ps.Size]
		if !ok || count == 0 {
			continue
		}
		packs = append(packs, Pack{
			Size:  ps.Size,
			Count: count,
		})
	}

	return packs, nil
}

func (s *Service) AddSize(ctx context.Context, size int64) (*PackSize, error) {
	created, err := s.repo.AddSize(ctx, size)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateSize) {
			return nil, ErrSizeExists
		}
		return nil, fmt.Errorf("failed to add pack size: %w", err)
	}

	return &PackSize{
		ID:   created.ID,
		Size: created.Size,
	}, nil
}

func (s *Service) DeleteSize(ctx context.Context, id int64) error {
	if err := s.repo.EnsureSizeExists(ctx, id); err != nil {
		if errors.Is(err, repository.ErrSizeNotFound) {
			return ErrSizeNotFound
		}
		return fmt.Errorf("failed to check pack size: %w", err)
	}

	count, err := s.repo.CountSizes(ctx)
	if err != nil {
		return fmt.Errorf("failed to count pack sizes: %w", err)
	}

	if count <= 1 {
		return ErrLastSize
	}

	if err := s.repo.DeleteSize(ctx, id); err != nil {
		if errors.Is(err, repository.ErrSizeNotFound) {
			return ErrSizeNotFound
		}
		return fmt.Errorf("failed to delete pack size: %w", err)
	}

	return nil
}
