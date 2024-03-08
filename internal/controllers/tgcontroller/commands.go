package tgcontroller

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
				Command:     "who",
				Description: "кто я, и где",
			},
			{
				Command:     "quote",
				Description: "великая цитата",
			},
			{
				Command:     "comment",
				Description: "сочный комментарий",
			},
			{
				Command:     "add_quote",
				Description: "добавить цитату",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("set commands: %w", err)
	}

	return nil
}
