package cmsusecases

import (
	"app/internal/domain"
	"context"
)

func (uc *UseCases) Moderators(ctx context.Context) ([]domain.Moderator, error) {
	return uc.repo.Moderators(ctx) // Пока просто проксируем
}

func (uc *UseCases) AddModerator(ctx context.Context, userID int64, description string) error {
	return uc.repo.AddModerator(ctx, userID, description) // Пока просто проксируем
}

func (uc *UseCases) DeleteModerator(ctx context.Context, userID int64) error {
	return uc.repo.DeleteModerator(ctx, userID) // Пока просто проксируем
}
