package tgcontroller

import (
	"app/internal/common"
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	subsystemName       = "controller"
	endpointNameLabel   = "endpoint"
	endpointStatusLabel = "status"
)

var (
	updateCount = promauto.With(common.DefaultRegistry).NewCounter(prometheus.CounterOpts{
		Namespace: common.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "update_count",
		Help:      "Количество обновлений (сообщений)",
	})
	handleCount = promauto.With(common.DefaultRegistry).NewCounterVec(prometheus.CounterOpts{
		Namespace: common.MetricsNamespace,
		Subsystem: subsystemName,
		Name:      "handle_count",
		Help:      "Количество обработанных сообщений",
	}, []string{endpointNameLabel, endpointStatusLabel})
)

func (c *Controller) counterMiddleware() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			updateCount.Inc()

			next(ctx, bot, update)
		}
	}
}

func (c *Controller) accessMiddleware() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			if update.Message != nil {
				for _, cid := range c.allowedChats {
					if update.Message.Chat.ID == cid {
						next(ctx, bot, update)

						return
					}
				}
			}

		}
	}
}

func (c *Controller) handleWrapper(next handler, name string) handler {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
		ok, err := next(ctx, b, update)

		if ok {
			handleCount.WithLabelValues(name, common.ConvertOk(err == nil)).Inc()
		}

		return ok, err
	}
}
