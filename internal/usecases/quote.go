package usecases

import (
	"context"
	"fmt"
)

func (u *UseCases) RandomQuote(ctx context.Context) (string, error) {
	quote, err := u.repo.RandomQuote(ctx)
	if err != nil {
		return "", fmt.Errorf("use case: random quote: %w", err)
	}

	return quote, nil
}
