package app

import (
	"context"
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

func (s *Service) Sizes(ctx context.Context) ([]int64, error) {
	return s.repo.Sizes(ctx)
}

func (s *Service) CalculatePackages(ctx context.Context, orderedItems int64) ([]Pack, error) {
	sizes, err := s.Sizes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sizes: %w", err)
	}

	result, err := calculate.CalculatePackages(orderedItems, sizes)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate packages: %w", err)
	}

	packs := make([]Pack, 0, len(result))
	for _, size := range sizes {
		count, ok := result[size]
		if !ok || count == 0 {
			continue
		}
		packs = append(packs, Pack{
			Size:  size,
			Count: count,
		})
	}

	return packs, nil
}
