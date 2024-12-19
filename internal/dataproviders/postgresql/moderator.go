package postgresql

import (
	"app/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) IsModerator(ctx context.Context, botID int64, userID int64) (bool, error) {
	raw := 0

	err := r.db.GetContext(ctx, &raw, `SELECT 1 FROM moderatos WHERE bot_id = $1 AND user_id = $2 LIMIT 1;`, botID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("select moderator: %w", err)
	}

	return true, nil
}

func (r *Repository) allModerators(ctx context.Context, botID int64) ([]*moderatorModel, error) {
	raw := make([]*moderatorModel, 0)

	err := r.db.SelectContext(ctx, &raw, `SELECT * FROM moderators WHERE bot_id = $1;`, botID)
	if err != nil {
		return nil, fmt.Errorf("all moderators: %w", err)
	}

	return raw, nil
}

func (r *Repository) Moderators(ctx context.Context, botID int64) ([]domain.Moderator, error) {
	rawModerators, err := r.allModerators(ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("repository: moderators: %w", err)
	}

	moderators := make([]domain.Moderator, len(rawModerators))
	for i, v := range rawModerators {
		moderators[i] = v.toDomain()
	}

	return moderators, nil
}

func (r *Repository) AddModerator(ctx context.Context, botID, userID int64, description string) error {
	builder := squirrel.Insert("moderators").
		PlaceholderFormat(squirrel.Dollar).
		SetMap(
			map[string]interface{}{
				"user_id":     userID,
				"bot_id":      botID,
				"description": StringToDB(description),
			},
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("repository: build query: %w", err)
	}

	r.squirrelDebugLog(ctx, query, args)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("repository: add moderator: %w", err)
	}

	return nil
}

func (r *Repository) DeleteModerator(ctx context.Context, botID, userID int64) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM moderators WHERE bot_id = $1 AND user_id = $2;`,
		botID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("repository: delete moderator: %w", err)
	}

	return nil
}
