package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Controller) whoHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}

	ok := strings.Index(update.Message.Text, "/who") == 0
	if !ok {
		return false, nil
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             fmt.Sprintf("user id: %d chat id: %d", update.Message.From.ID, update.Message.Chat.ID),
		ReplyToMessageID: update.Message.ID,
	})
	if err != nil {
		return true, fmt.Errorf("who handle: send message: %w", err)
	}

	return true, nil
}
