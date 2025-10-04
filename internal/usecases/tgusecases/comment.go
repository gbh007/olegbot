package tgusecases

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *UseCases) CommentHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil || update.Message.ReplyToMessage == nil {
		return false, nil
	}

	ok, err := u.commandStrictCheck(ctx, "/comment", update.Message.Text)
	if err != nil {
		return true, fmt.Errorf("comment handle: strict check: %w", err)
	}

	if !ok {
		return false, nil
	}

	quote, err := u.randomQuote(ctx, nil, []string{update.Message.Text}, false)
	if err != nil {
		return true, fmt.Errorf("comment handle: %w", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   quote,
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ReplyToMessage.ID,
			ChatID:    update.Message.Chat.ID,
		},
	})
	if err != nil {
		return true, fmt.Errorf("comment handle: send message: %w", err)
	}

	_, _ = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})

	return true, nil
}
