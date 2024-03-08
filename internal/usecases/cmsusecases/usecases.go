package cmsusecases

import (
	"app/internal/domain"
	"context"
)

type repository interface {
	AddQuote(ctx context.Context, text string, userID, chatID int64) error
	QuoteExists(ctx context.Context, text string) (bool, error)
	Quotes(ctx context.Context) ([]domain.Quote, error)
}

type UseCases struct {
	repo repository
}

func New(repo repository) *UseCases {
	return &UseCases{
		repo: repo,
	}
}
