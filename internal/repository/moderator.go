package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type moderatorModel struct {
	UserID      int64          `db:"user_id"`
	CreateAt    time.Time      `db:"create_at"`
	Description sql.NullString `db:"description"`
}

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
