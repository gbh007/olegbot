package postgresql

import (
	"app/internal/domain"
	"context"
	"fmt"
)

func (r *Repository) allQuotes(ctx context.Context, botID int64) ([]*quoteModel, error) {
	raw := make([]*quoteModel, 0)

	err := r.db.SelectContext(ctx, &raw, `SELECT * FROM "quotes" WHERE bot_id = $1;`, botID)
	if err != nil {
		return nil, fmt.Errorf("all quotes: %w", err)
	}

	return raw, nil
}

func (r *Repository) Quotes(ctx context.Context, botID int64) ([]domain.Quote, error) {
	rawQuotes, err := r.allQuotes(ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("repository: quotes: %w", err)
	}

	quotes := make([]domain.Quote, len(rawQuotes))
	for i, v := range rawQuotes {
		quotes[i] = v.toDomain()
	}

	return quotes, nil
}
