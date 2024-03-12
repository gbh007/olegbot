package cmsusecases

import (
	"app/internal/domain"
	"context"
	"fmt"
)

func (uc *UseCases) Quote(ctx context.Context, id int64) (domain.Quote, error) {
	return uc.repo.Quote(ctx, id) // Пока просто проксируем
}

func (uc *UseCases) DeleteQuote(ctx context.Context, id int64) error {
	return uc.repo.DeleteQuote(ctx, id) // Пока просто проксируем
}

func (uc *UseCases) UpdateQuoteText(ctx context.Context, id int64, text string) error {
	return uc.repo.UpdateQuoteText(ctx, id, text) // Пока просто проксируем
}

func (u *UseCases) AddQuote(ctx context.Context, text string) error {
	exists, err := u.repo.QuoteExists(ctx, text)
	if err != nil {
		return fmt.Errorf("use case: add quote: %w", err)
	}

	if exists {
		return fmt.Errorf("use case: add quote: %w", domain.QuoteAlreadyExistsError)
	}

	err = u.repo.AddQuote(ctx, text, 0, 0)
	if err != nil {
		return fmt.Errorf("use case: add quote: %w", err)
	}

	return nil
}

func (u *UseCases) AddQuotes(ctx context.Context, quotes []string) error {
	for _, text := range quotes {
		exists, err := u.repo.QuoteExists(ctx, text)
		if err != nil {
			return fmt.Errorf("use case: add quotes: %w", err)
		}

		if exists {
			continue
		}

		err = u.repo.AddQuote(ctx, text, 0, 0)
		if err != nil {
			return fmt.Errorf("use case: add quotes: %w", err)
		}
	}

	return nil
}
