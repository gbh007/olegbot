package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type quoteModel struct {
	ID       int64         `db:"id"`
	Text     string        `db:"text"`
	CreateAt time.Time     `db:"create_at"`
	UserID   sql.NullInt64 `db:"user_id"`
	ChatID   sql.NullInt64 `db:"chat_id"`
}

func (r *Repository) RandomQuote(_ context.Context) (string, error) {
	r.dataMutex.RLock()
	defer r.dataMutex.RUnlock()

	return r.data[rand.Intn(len(r.data))], nil
}

func (r *Repository) AddQuote(ctx context.Context, text string, userID, chatID int64) error {
	r.dataMutex.Lock()
	defer r.dataMutex.Unlock()

	_, err := r.db.Exec(
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

	r.data = append(r.data, text)

	return nil
}

func (r *Repository) allQuotes(ctx context.Context) ([]*quoteModel, error) {
	raw := make([]*quoteModel, 0)

	err := r.db.SelectContext(ctx, &raw, `SELECT * FROM "quotes";`)
	if err != nil {
		return nil, fmt.Errorf("all quotes: %w", err)
	}

	return raw, nil
}
