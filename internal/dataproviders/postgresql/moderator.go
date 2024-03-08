package postgresql

import (
	"app/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

func (r *Repository) IsModerator(_ context.Context, userID int64) (bool, error) {
	r.moderatorsMutex.RLock()
	defer r.moderatorsMutex.RUnlock()

	_, ok := r.moderators[userID]

	return ok, nil
}

func (r *Repository) allModerators(ctx context.Context) ([]*moderatorModel, error) {
	raw := make([]*moderatorModel, 0)

	err := r.db.SelectContext(ctx, &raw, `SELECT * FROM moderators;`)
	if err != nil {
		return nil, fmt.Errorf("all moderators: %w", err)
	}

	return raw, nil
}

func (r *Repository) Moderators(ctx context.Context) ([]domain.Moderator, error) {
	rawModerators, err := r.allModerators(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository: moderators: %w", err)
	}

	moderators := make([]domain.Moderator, len(rawModerators))
	for i, v := range rawModerators {
		moderators[i] = v.toDomain()
	}

	return moderators, nil
}

func (r *Repository) AddModerator(ctx context.Context, userID int64, description string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO moderators (user_id, "description") VALUES ($1, $2);`,
		userID,
		sql.NullString{
			String: description,
			Valid:  description != "",
		},
	)
	if err != nil {
		return fmt.Errorf("repository: add moderator: %w", err)
	}

	err = r.reloadModeratorCache(ctx)
	if err != nil {
		return fmt.Errorf("repository: add moderator: %w", err)
	}

	return nil
}

func (r *Repository) DeleteModerator(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM moderators WHERE user_id = $1;`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("repository: delete moderator: %w", err)
	}

	err = r.reloadModeratorCache(ctx)
	if err != nil {
		return fmt.Errorf("repository: delete moderator: %w", err)
	}

	return nil
}
