package postgresql

import (
	"app/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
)

func (r *Repository) RandomQuote(_ context.Context) (string, error) {
	r.dataMutex.RLock()
	defer r.dataMutex.RUnlock()

	return r.data[rand.Intn(len(r.data))], nil
}

func (r *Repository) AddQuote(ctx context.Context, text string, userID, chatID int64) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO "quotes" ("text", user_id, chat_id) VALUES ($1, $2, $3);`,
		text,
		sql.NullInt64{
			Int64: userID,
			Valid: userID != 0,
		},
		sql.NullInt64{
			Int64: chatID,
			Valid: chatID != 0,
		},
	)
	if err != nil {
		return fmt.Errorf("repository: add quote: %w", err)
	}

	err = r.reloadQuoteCache(ctx)
	if err != nil {
		return fmt.Errorf("repository: add quote: %w", err)
	}

	return nil
}

func (r *Repository) QuoteExists(ctx context.Context, text string) (bool, error) {
	var count int64

	err := r.db.GetContext(
		ctx,
		&count,
		`SELECT COUNT(*) FROM "quotes" WHERE LOWER("text") = $1;`,
		strings.ToLower(text),
	)
	if err != nil {
		return false, fmt.Errorf("repository: quote exists: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) Quote(ctx context.Context, id int64) (domain.Quote, error) {
	var rawQuote quoteModel

	err := r.db.GetContext(
		ctx,
		&rawQuote,
		`SELECT * FROM "quotes" WHERE id = $1 LIMIT 1;`,
		id,
	)
	if err != nil {
		return domain.Quote{}, fmt.Errorf("repository: quote: %w", err)
	}

	return rawQuote.toDomain(), nil
}

func (r *Repository) DeleteQuote(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM "quotes" WHERE id = $1;`,
		id,
	)
	if err != nil {
		return fmt.Errorf("repository: delete quote: %w", err)
	}

	err = r.reloadQuoteCache(ctx)
	if err != nil {
		return fmt.Errorf("repository: delete quote: %w", err)
	}

	return nil
}

func (r *Repository) UpdateQuoteText(ctx context.Context, id int64, text string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE "quotes" SET text = $2 WHERE id = $1;`,
		id, text,
	)
	if err != nil {
		return fmt.Errorf("repository: update quote text: %w", err)
	}

	err = r.reloadQuoteCache(ctx)
	if err != nil {
		return fmt.Errorf("repository: update quote text: %w", err)
	}

	return nil
}
