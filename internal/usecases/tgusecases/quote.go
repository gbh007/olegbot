package tgusecases

import (
	"app/internal/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *UseCases) RandomQuote(ctx context.Context) (string, error) {
	quote, err := u.repo.RandomQuote(ctx)
	if err != nil {
		return "", fmt.Errorf("use case: random quote: %w", err)
	}

	return quote, nil
}

func (u *UseCases) AddQuote(ctx context.Context, text string, userID, chatID int64) error {
	ok, err := u.repo.IsModerator(ctx, userID)
	if err != nil {
		return fmt.Errorf("use case: add quote: %w", err)
	}

	if !ok {
		return fmt.Errorf("use case: add quote: %w", domain.PermissionDeniedError)
	}

	exists, err := u.repo.QuoteExists(ctx, text)
	if err != nil {
		return fmt.Errorf("use case: add quote: %w", err)
	}

	if exists {
		return fmt.Errorf("use case: add quote: %w", domain.QuoteAlreadyExistsError)
	}

	err = u.repo.AddQuote(ctx, text, userID, chatID)
	if err != nil {
		return fmt.Errorf("use case: add quote: %w", err)
	}

	return nil
}

func (u *UseCases) QuoteHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}

	ok := strings.Index(update.Message.Text, "/quote") == 0
	if !ok {
		return false, nil
	}

	quote, err := u.RandomQuote(ctx)
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
		ChatID: update.Message.Chat.ID,
		Text:   quote,
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMessageID,
			ChatID:    update.Message.Chat.ID,
		},
	})
	if err != nil {
		return true, fmt.Errorf("quote handle: send message: %w", err)
	}

	return true, nil
}

func (u *UseCases) AddQuoteHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
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

	var (
		text    string
		replyTo int
	)

	switch {
	case update.Message.ReplyToMessage != nil:
		text = update.Message.ReplyToMessage.Text
		replyTo = update.Message.ReplyToMessage.ID
	default:
		text = update.Message.Text[splitIndex+1:]
		replyTo = update.Message.ID
	}

	err := u.AddQuote(ctx, text, update.Message.From.ID, update.Message.Chat.ID)
	switch {
	case errors.Is(err, domain.QuoteAlreadyExistsError):
		_, sendErr := b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: replyTo,
			Reaction: []models.ReactionType{
				{
					Type: models.ReactionTypeTypeEmoji,
					ReactionTypeEmoji: &models.ReactionTypeEmoji{
						Type:  models.ReactionTypeTypeEmoji,
						Emoji: "👎", // FIXME: в конфигурацию
					},
				},
			},
		})
		if sendErr != nil {
			return true, fmt.Errorf("add quote handle: send message: %w", sendErr)
		}

	case err != nil:
		return true, fmt.Errorf("add quote handle: %w", err)

	default:
		_, sendErr := b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: replyTo,
			Reaction: []models.ReactionType{
				{
					Type: models.ReactionTypeTypeEmoji,
					ReactionTypeEmoji: &models.ReactionTypeEmoji{
						Type:  models.ReactionTypeTypeEmoji,
						Emoji: "👍", // FIXME: в конфигурацию
					},
				},
			},
		})
		if sendErr != nil {
			return true, fmt.Errorf("add quote handle: send message: %w", sendErr)
		}
	}

	if update.Message.ReplyToMessage != nil {
		_, deleteErr := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		})
		if deleteErr != nil {
			return true, fmt.Errorf("add quote handle: delete message: %w", deleteErr)
		}
	}

	return true, nil
}
