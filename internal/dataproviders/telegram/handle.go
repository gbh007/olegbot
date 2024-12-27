package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Controller) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	defer func() {
		_ = recover()
	}()

	for _, h := range c.handlers {
		// FIXME: обрабатывать ошибку
		ok, _ := h(ctx, b, update)
		if ok {
			break
		}
	}
}
