package tgusecases

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *UseCases) AccessMiddleware() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			if update.Message != nil {
				botInfo, err := u.repo.BotInfo(ctx)
				if err != nil {
					// TODO: логировать ошибку
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
