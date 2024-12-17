package tgusecases

import (
	"app/internal/domain"
	"context"
)

type repository interface {
	RandomQuote(ctx context.Context) (string, error)
	AddQuote(ctx context.Context, text string, userID, chatID int64) error
	IsModerator(ctx context.Context, userID int64) (bool, error)
	QuoteExists(ctx context.Context, text string) (bool, error)
	BotInfo(ctx context.Context) (domain.Bot, error)
}

type UseCases struct {
	repo repository
}

func New(
	repo repository,
	emojiList []string,
	emojiChance float32,
	rawTags []string,
	botName, botTag string,
) *UseCases {
	return &UseCases{
		repo: repo,
	}
}
