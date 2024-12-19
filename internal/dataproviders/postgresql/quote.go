package postgresql

import (
	"app/internal/domain"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) AddQuote(ctx context.Context, botID int64, text string, userID, chatID int64) error {
	builder := squirrel.Insert("quotes").
		PlaceholderFormat(squirrel.Dollar).
		SetMap(
			map[string]interface{}{
				"bot_id":  botID,
				"text":    text,
				"user_id": Int64ToDB(userID),
				"chat_id": Int64ToDB(chatID),
			},
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("repository: build query: %w", err)
	}

	r.squirrelDebugLog(ctx, query, args)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("repository: add quote: %w", err)
	}

	return nil
}

func (r *Repository) QuoteExists(ctx context.Context, botID int64, text string) (bool, error) {
	var count int64

	err := r.db.GetContext(
		ctx,
		&count,
		`SELECT COUNT(*) FROM "quotes" WHERE bot_id = $1 AND LOWER("text") = $2;`,
		botID, strings.ToLower(text),
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

	return nil
}

func (r *Repository) UpdateQuoteText(ctx context.Context, id int64, text string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE "quotes" SET text = $2, update_at = $3 WHERE id = $1;`,
		id, text, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("repository: update quote text: %w", err)
	}

	return nil
}
