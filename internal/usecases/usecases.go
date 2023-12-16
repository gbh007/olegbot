package usecases

import "context"

type repository interface {
	RandomQuote(ctx context.Context) (string, error)
	AddQuote(ctx context.Context, text string, userID, chatID int64) error
	IsModerator(ctx context.Context, userID int64) (bool, error)
	QuoteExists(ctx context.Context, text string) (bool, error)
}

type UseCases struct {
	repo repository
}

func New(repo repository) *UseCases {
	return &UseCases{
		repo: repo,
	}
}
