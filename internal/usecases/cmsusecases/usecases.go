package cmsusecases

import (
	"app/internal/domain"
	"context"
)

type repository interface {
	AddQuote(ctx context.Context, botID int64, text string, userID, chatID int64) error
	UpdateQuoteText(ctx context.Context, id int64, text string) error
	DeleteQuote(ctx context.Context, id int64) error

	QuoteExists(ctx context.Context, botID int64, text string) (bool, error)

	Quotes(ctx context.Context, botID int64) ([]domain.Quote, error)
	Quote(ctx context.Context, id int64) (domain.Quote, error)

	Moderators(ctx context.Context, botID int64) ([]domain.Moderator, error)
	AddModerator(ctx context.Context, botID, userID int64, description string) error
	DeleteModerator(ctx context.Context, botID, userID int64) error
}

type UseCases struct {
	repo repository
}

func New(repo repository) *UseCases {
	return &UseCases{
		repo: repo,
	}
}
