package telegram

import (
	"app/internal/common"
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Controller) counterMiddleware() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			common.UpdateCount.Inc()

			next(ctx, bot, update)
		}
	}
}

func (c *Controller) handleWrapper(next handler, name string) handler {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
		ok, err := next(ctx, b, update)

		if ok {
			common.HandleCount.WithLabelValues(name, common.ConvertOk(err == nil)).Inc()
		}

		return ok, err
	}
}
