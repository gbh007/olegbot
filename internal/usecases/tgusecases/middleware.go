package tgusecases

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *UseCases) AccessMiddleware() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			if update.Message != nil {
				botInfo, err := u.repo.GetBot(ctx, u.botID)
				if err != nil {
					if u.debug {
						u.logger.DebugContext(ctx, "access middleware check", slog.String("error", err.Error()))
					}

					return
				}

				if len(botInfo.AllowedChats) == 0 {
					next(ctx, bot, update)

					return
				}

				for _, cid := range botInfo.AllowedChats {
					if update.Message.Chat.ID == cid {
						next(ctx, bot, update)

						return
					}
				}
			}
		}
	}
}
