package usecases

import (
	"app/internal/domain"
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
		return fmt.Errorf("use case: add quote: %w", domain.PermissionDeniedError)
	}

	exists, err := u.repo.QuoteExists(ctx, text)
	if err != nil {
		return fmt.Errorf("use case: add quote: %w", err)
	}

	if exists {
		return fmt.Errorf("use case: add quote: %w", domain.QuoteAlreadyExistsError)
	}

	err = u.repo.AddQuote(ctx, text, userID, chatID)
	if err != nil {
		return fmt.Errorf("use case: add quote: %w", err)
	}

	return nil
}
