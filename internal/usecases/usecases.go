package usecases

import "context"

type repository interface {
	AddQuote(ctx context.Context, text string) error
	RandomQuote(ctx context.Context) (string, error)
}

type UseCases struct {
	repo repository
}

func New(repo repository) *UseCases {
	return &UseCases{
		repo: repo,
	}
}
