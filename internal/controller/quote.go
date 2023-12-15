package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

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

func (c *Controller) addQuoteHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}

	ok := strings.Index(update.Message.Text, "/add_quote") == 0
	if !ok {
		return false, nil
	}

	splitIndex := strings.Index(update.Message.Text, " ")
	if splitIndex == -1 && update.Message.ReplyToMessage == nil {
		return true, fmt.Errorf("add quote invalid syntax")
	}

	var text string

	switch {
	case update.Message.ReplyToMessage != nil:
		text = update.Message.ReplyToMessage.Text
	default:
		text = update.Message.Text[splitIndex+1:]
	}

	err := c.useCases.AddQuote(ctx, text, update.Message.From.ID, update.Message.Chat.ID)
	if err != nil {
		return true, fmt.Errorf("add quote handle: %w", err)
	}

	_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})
	if err != nil {
		return true, fmt.Errorf("add quote handle: delete message: %w", err)
	}

	return true, nil
}
