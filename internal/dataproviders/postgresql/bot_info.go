package postgresql

import (
	"app/internal/domain"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) BotInfo(ctx context.Context, botID int64) (domain.Bot, error) {
	raw := botModel{}

	builder := squirrel.Select(raw.columns()...).
		PlaceholderFormat(squirrel.Dollar).
		From("bots").
		Where(squirrel.Eq{"id": botID}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return domain.Bot{}, fmt.Errorf("repository: select bot: build query: %w", err)
	}

	r.squirrelDebugLog(ctx, query, args)

	err = r.pool.QueryRow(ctx, query, args...).Scan(&raw)
	if err != nil {
		return domain.Bot{}, fmt.Errorf("repository: select bot: scan: %w", err)
	}

	return raw.toDomain(), nil
}
