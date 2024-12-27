package tgusecases

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *UseCases) WhoHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}

	ok, err := u.commandStrictCheck(ctx, "/who", update.Message.Text)
	if err != nil {
		return true, fmt.Errorf("who handle: strict check: %w", err)
	}

	if !ok {
		return false, nil
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("user id: %d chat id: %d", update.Message.From.ID, update.Message.Chat.ID),
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
			ChatID:    update.Message.Chat.ID,
		},
	})
	if err != nil {
		return true, fmt.Errorf("who handle: send message: %w", err)
	}

	return true, nil
}
