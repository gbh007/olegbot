package postgresql

import (
	"app/internal/domain"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) GetBot(ctx context.Context, botID int64) (domain.Bot, error) {
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

func (r *Repository) CreateBot(ctx context.Context, bot domain.Bot) error {
	raw := botModel{}
	raw.fromDomain(bot)

	builder := squirrel.Insert("bots").
		PlaceholderFormat(squirrel.Dollar).
		SetMap(map[string]interface{}{
			"name":           raw.Name,
			"bot_tag":        raw.BotTag,
			"token":          raw.Token,
			"enabled":        raw.Enabled,
			"description":    raw.Description,
			"emoji_list":     raw.EmojiList,
			"emoji_chance":   raw.EmojiChance,
			"tags":           raw.Tags,
			"allowed_chats":  raw.AllowedChats,
			"sticker_chance": raw.StickerChance,
			"gif_chance":     raw.GifChance,
			"create_at":      raw.CreateAt,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("repository: insert bot: build query: %w", err)
	}

	r.squirrelDebugLog(ctx, query, args)

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("repository: insert bot: exec query: %w", err)
	}

	return nil
}

func (r *Repository) UpdateBot(ctx context.Context, bot domain.Bot) error {
	raw := botModel{}
	raw.fromDomain(bot)

	builder := squirrel.Update("bots").
		PlaceholderFormat(squirrel.Dollar).
		SetMap(map[string]interface{}{
			"name":           raw.Name,
			"bot_tag":        raw.BotTag,
			"token":          raw.Token,
			"enabled":        raw.Enabled,
			"description":    raw.Description,
			"emoji_list":     raw.EmojiList,
			"emoji_chance":   raw.EmojiChance,
			"tags":           raw.Tags,
			"allowed_chats":  raw.AllowedChats,
			"sticker_chance": raw.StickerChance,
			"gif_chance":     raw.GifChance,
			"update_at":      raw.UpdateAt,
		}).
		Where(squirrel.Eq{
			"id": raw.ID,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("repository: update bot: build query: %w", err)
	}

	r.squirrelDebugLog(ctx, query, args)

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("repository: update bot: exec query: %w", err)
	}

	return nil
}

func (r *Repository) DeleteBot(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM bots WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("repository: delete bot: exec query: %w", err)
	}

	return nil
}

func (r *Repository) GetBots(ctx context.Context) ([]domain.Bot, error) {
	builder := squirrel.Select(botModel{}.columns()...).
		PlaceholderFormat(squirrel.Dollar).
		From("bots")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("repository: select bot: build query: %w", err)
	}

	r.squirrelDebugLog(ctx, query, args)

	result := make([]domain.Bot, 0)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: select bot: exec: %w", err)
	}

	for rows.Next() {
		raw := botModel{}

		err = rows.Scan(&raw)
		if err != nil {
			return nil, fmt.Errorf("repository: select bot: scan: %w", err)
		}

		result = append(result, raw.toDomain())
	}

	return result, nil
}
