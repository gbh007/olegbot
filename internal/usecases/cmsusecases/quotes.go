package cmsusecases

import (
	"app/internal/domain"
	"context"
)

func (uc *UseCases) Quotes(ctx context.Context, botID int64) ([]domain.Quote, error) {
	return uc.repo.Quotes(ctx, botID) // Пока просто проксируем
}
