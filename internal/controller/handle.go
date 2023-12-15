package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Controller) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	for _, h := range c.handlers {
		// FIXME: обрабатывать ошибку
		ok, _ := h(ctx, b, update)
		if ok {
			break
		}
	}
}

func (c *Controller) commentHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil || update.Message.ReplyToMessage == nil {
		return false, nil
	}

	ok := strings.Index(update.Message.Text, "/comment") == 0
	if !ok {
		return false, nil
	}

	quote, err := c.useCases.RandomQuote(ctx)
	if err != nil {
		return true, fmt.Errorf("comment handle: %w", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             quote,
		ReplyToMessageID: update.Message.ReplyToMessage.ID,
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

func (c *Controller) quoteHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}

	ok := strings.Index(update.Message.Text, "/quote") == 0
	if !ok {
		return false, nil
	}

	quote, err := c.useCases.RandomQuote(ctx)
	if err != nil {
		return true, fmt.Errorf("quote handle: %w", err)
	}

	replyToMessageID := update.Message.ID

	ok, _ = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})

	if ok {
		replyToMessageID = 0
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             quote,
		ReplyToMessageID: replyToMessageID,
	})
	if err != nil {
		return true, fmt.Errorf("quote handle: send message: %w", err)
	}

	return true, nil
}

func (c *Controller) selfHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}

	messageText := strings.ToLower(update.Message.Text)

	ok := c.hasBotName &&strings.Contains(messageText, c.botName) ||
		c.hasBotTag && strings.Contains(messageText, c.botTag)
	if !ok {
		return false, nil
	}

	quote, err := c.useCases.RandomQuote(ctx)
	if err != nil {
		return true, fmt.Errorf("self handle: %w", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             quote,
		ReplyToMessageID: update.Message.ID,
	})
	if err != nil {
		return true, fmt.Errorf("self handle: send message: %w", err)
	}

	return true, nil
}
