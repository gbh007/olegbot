package telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Controller) setBotCommands(ctx context.Context, b *bot.Bot) error {
	_, err := b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{
				Command:     "comment",
				Description: "сочный комментарий",
			},
			{
				Command:     "sticker",
				Description: "крутой стикер",
			},
			{
				Command:     "gif",
				Description: "горячая гифка",
			},
			{
				Command:     "emoji",
				Description: "топовая реакция",
			},
			{
				Command:     "quote",
				Description: "великая цитата",
			},
			{
				Command:     "who",
				Description: "кто я, и где",
			},
			{
				Command:     "add_quote",
				Description: "добавить цитату",
			},
			{
				Command:     "add_sticker",
				Description: "добавить стикер",
			},
			{
				Command:     "add_gif",
				Description: "добавить гифку",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("set commands: %w", err)
	}

	return nil
}
