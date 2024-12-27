package cmsusecases

import (
	"app/internal/domain"
	"context"
)

func (uc *UseCases) Moderators(ctx context.Context, botID int64) ([]domain.Moderator, error) {
	return uc.repo.Moderators(ctx, botID) // Пока просто проксируем
}

func (uc *UseCases) AddModerator(ctx context.Context, botID, userID int64, description string) error {
	return uc.repo.AddModerator(ctx, botID, userID, description) // Пока просто проксируем
}

func (uc *UseCases) DeleteModerator(ctx context.Context, botID, userID int64) error {
	return uc.repo.DeleteModerator(ctx, botID, userID) // Пока просто проксируем
}
