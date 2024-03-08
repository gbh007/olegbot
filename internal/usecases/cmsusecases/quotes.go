package cmsusecases

import (
	"app/internal/domain"
	"context"
)

func (uc *UseCases) Quotes(ctx context.Context) ([]domain.Quote, error) {
	return uc.repo.Quotes(ctx) // Пока просто проксируем
}
