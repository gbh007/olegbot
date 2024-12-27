package telegram

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Controller) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	defer func() {
		_ = recover()
	}()

	for _, h := range c.handlers {
		ok, err := h(ctx, b, update)

		if c.debug && err != nil {
			c.logger.DebugContext(ctx, "handle tg update error", slog.String("error", err.Error()))
		}

		if ok {
			break
		}
	}
}
