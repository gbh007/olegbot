package postgresql

import (
	"app/internal/domain"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) AddGif(ctx context.Context, gif domain.Gif) error {
	builder := squirrel.Insert("gifs").
		PlaceholderFormat(squirrel.Dollar).
		SetMap(
			map[string]interface{}{
				"bot_id":  gif.BotID,
				"file_id": gif.FileID,
				"user_id": Int64ToDB(gif.CreatorID),
				"chat_id": Int64ToDB(gif.CreatedInChatID),
			},
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("repository: build query: %w", err)
	}

	r.squirrelDebugLog(ctx, query, args)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("repository: add gif: %w", err)
	}

	return nil
}

func (r *Repository) GifExists(ctx context.Context, botID int64, fileID string) (bool, error) {
	var count int64

	err := r.db.GetContext(
		ctx,
		&count,
		`SELECT COUNT(*) FROM gifs WHERE bot_id = $1 AND file_id = $2;`,
		botID, fileID,
	)
	if err != nil {
		return false, fmt.Errorf("repository: gif exists: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) Gifs(ctx context.Context, botID int64) ([]domain.Gif, error) {
	raw := make([]*gifModel, 0)

	err := r.db.SelectContext(ctx, &raw, `SELECT * FROM gifs WHERE bot_id = $1;`, botID)
	if err != nil {
		return nil, fmt.Errorf("repository: gifs: %w", err)
	}

	gifs := make([]domain.Gif, len(raw))
	for i, v := range raw {
		gifs[i] = v.toDomain()
	}

	return gifs, nil
}
