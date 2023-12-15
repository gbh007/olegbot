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

func (u *UseCases) AddQuote(ctx context.Context, text string, userID, chatID int64) error {
	ok, err := u.repo.IsModerator(ctx, userID)
	if err != nil {
		return fmt.Errorf("use case: add quote: %w", err)
	}

	if !ok {
		return fmt.Errorf("use case: add quote: not a moderator")
	}

	err = u.repo.AddQuote(ctx, text, userID, chatID)
	if err != nil {
		return fmt.Errorf("use case: add quote: %w", err)
	}

	return nil
}
