package cmsusecases

import (
	"app/internal/domain"
	"context"
	"time"
)

func (u *UseCases) CreateBot(ctx context.Context, bot domain.Bot) error {
	bot.CreateAt = time.Now().UTC()

	return u.repo.CreateBot(ctx, bot)
}
func (u *UseCases) UpdateBot(ctx context.Context, bot domain.Bot) error {
	bot.UpdateAt = time.Now().UTC()

	return u.repo.UpdateBot(ctx, bot)
}
func (u *UseCases) DeleteBot(ctx context.Context, id int64) error {
	return u.repo.DeleteBot(ctx, id)
}
func (u *UseCases) GetBots(ctx context.Context) ([]domain.Bot, error) {
	return u.repo.GetBots(ctx)
}
