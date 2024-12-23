package tgusecases

import (
	"app/internal/domain"
	"context"
)

type repository interface {
	Quotes(ctx context.Context, botID int64) ([]domain.Quote, error)
	AddQuote(ctx context.Context, botID int64, text string, userID, chatID int64) error
	IsModerator(ctx context.Context, botID int64, userID int64) (bool, error)
	QuoteExists(ctx context.Context, botID int64, text string) (bool, error)
	GetBot(ctx context.Context, botID int64) (domain.Bot, error)
}

type UseCases struct {
	repo  repository
	botID int64
}

func New(repo repository, botID int64) *UseCases {
	return &UseCases{
		repo:  repo,
		botID: botID,
	}
}
