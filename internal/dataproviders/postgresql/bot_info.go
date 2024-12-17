package postgresql

import (
	"app/internal/domain"
	"context"
)

func (r *Repository) BotInfo(ctx context.Context) (domain.Bot, error) {
	return r.botInfo, nil
}
