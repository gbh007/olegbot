package postgresql

import (
	"app/internal/domain"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) AddSticker(ctx context.Context, sticker domain.Sticker) error {
	builder := squirrel.Insert("stickers").
		PlaceholderFormat(squirrel.Dollar).
		SetMap(
			map[string]interface{}{
				"bot_id":  sticker.BotID,
				"file_id": sticker.FileID,
				"user_id": Int64ToDB(sticker.CreatorID),
				"chat_id": Int64ToDB(sticker.CreatedInChatID),
			},
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("repository: build query: %w", err)
	}

	r.squirrelDebugLog(ctx, query, args)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("repository: add sticker: %w", err)
	}

	return nil
}

func (r *Repository) StickerExists(ctx context.Context, botID int64, fileID string) (bool, error) {
	var count int64

	err := r.db.GetContext(
		ctx,
		&count,
		`SELECT COUNT(*) FROM stickers WHERE bot_id = $1 AND file_id = $2;`,
		botID, fileID,
	)
	if err != nil {
		return false, fmt.Errorf("repository: sticker exists: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) Stickers(ctx context.Context, botID int64) ([]domain.Sticker, error) {
	raw := make([]*stickerModel, 0)

	err := r.db.SelectContext(ctx, &raw, `SELECT * FROM stickers WHERE bot_id = $1;`, botID)
	if err != nil {
		return nil, fmt.Errorf("repository: stickers: %w", err)
	}

	stickers := make([]domain.Sticker, len(raw))
	for i, v := range raw {
		stickers[i] = v.toDomain()
	}

	return stickers, nil
}
