package tgusecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *UseCases) SelfHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}

	messageText := strings.ToLower(update.Message.Text)
	captionText := strings.ToLower(update.Message.Caption)

	tags, err := u.repo.Tags(ctx)
	if err != nil {
		return false, fmt.Errorf("self handle: get tags: %w", err)
	}

	if len(tags) == 0 {
		return false, nil
	}

	found := false
	for _, tag := range tags {
		if strings.Contains(messageText, tag) || strings.Contains(captionText, tag) {
			found = true
			break
		}
	}

	if !found {
		return false, nil
	}

	quote, err := u.RandomQuote(ctx)
	if err != nil {
		return true, fmt.Errorf("self handle: %w", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   quote,
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
			ChatID:    update.Message.Chat.ID,
		},
	})
	if err != nil {
		return true, fmt.Errorf("self handle: send message: %w", err)
	}

	return true, nil
}
