package tgusecases

import (
	"app/internal/domain"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *UseCases) randomGif(ctx context.Context) (string, error) {
	gifs, err := u.repo.Gifs(ctx, u.botID)
	if err != nil {
		return "", fmt.Errorf("use case: random gif: %w", err)
	}

	if len(gifs) == 0 {
		return "", fmt.Errorf("use case: random gif: no gifs")
	}

	return gifs[rand.Intn(len(gifs))].FileID, nil
}

func (u *UseCases) addGif(ctx context.Context, fileID string, userID, chatID int64) error {
	ok, err := u.repo.IsModerator(ctx, u.botID, userID)
	if err != nil {
		return fmt.Errorf("use case: add gif: %w", err)
	}

	if !ok {
		return fmt.Errorf("use case: add gif: %w", domain.PermissionDeniedError)
	}

	exists, err := u.repo.GifExists(ctx, u.botID, fileID)
	if err != nil {
		return fmt.Errorf("use case: add gif: %w", err)
	}

	if exists {
		return fmt.Errorf("use case: add gif: %w", domain.AlreadyExistsError)
	}

	err = u.repo.AddGif(ctx, domain.Gif{
		BotID:           u.botID,
		FileID:          fileID,
		CreatorID:       userID,
		CreatedInChatID: chatID,
	})
	if err != nil {
		return fmt.Errorf("use case: add gif: %w", err)
	}

	return nil
}

func (u *UseCases) GifHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil || update.Message.ReplyToMessage == nil {
		return false, nil
	}

	ok, err := u.commandStrictCheck(ctx, "/gif", update.Message.Text)
	if err != nil {
		return true, fmt.Errorf("gif handle: strict check: %w", err)
	}

	if !ok {
		return false, nil
	}

	fileID, err := u.randomGif(ctx)
	if err != nil {
		return true, fmt.Errorf("gif handle: %w", err)
	}

	_, err = b.SendAnimation(ctx, &bot.SendAnimationParams{
		ChatID: update.Message.Chat.ID,
		Animation: &models.InputFileString{
			Data: fileID,
		},
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ReplyToMessage.ID,
			ChatID:    update.Message.Chat.ID,
		},
	})
	if err != nil {
		return true, fmt.Errorf("gif handle: send message: %w", err)
	}

	_, _ = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})

	return true, nil
}

func (u *UseCases) AddGifHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}

	ok, err := u.commandStrictCheck(ctx, "/add_gif", update.Message.Text)
	if err != nil {
		return true, fmt.Errorf("add gif handle: strict check: %w", err)
	}

	if !ok {
		return false, nil
	}

	splitIndex := strings.Index(update.Message.Text, " ")
	if splitIndex == -1 && (update.Message.ReplyToMessage == nil || update.Message.ReplyToMessage.Animation == nil) {
		return true, fmt.Errorf("add gif invalid syntax")
	}

	var (
		fileID  string
		replyTo int
	)

	switch {
	case update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.Animation != nil:
		fileID = update.Message.ReplyToMessage.Animation.FileID
		replyTo = update.Message.ReplyToMessage.ID
	default:
		fileID = update.Message.Text[splitIndex+1:]
		replyTo = update.Message.ID
	}

	err = u.addGif(ctx, fileID, update.Message.From.ID, update.Message.Chat.ID)
	switch {
	case errors.Is(err, domain.AlreadyExistsError):
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
			return true, fmt.Errorf("add gif handle: send message: %w", sendErr)
		}

	case err != nil:
		return true, fmt.Errorf("add gif handle: %w", err)

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
			return true, fmt.Errorf("add gif handle: send message: %w", sendErr)
		}
	}

	if update.Message.ReplyToMessage != nil {
		_, deleteErr := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		})
		if deleteErr != nil {
			return true, fmt.Errorf("add gif handle: delete message: %w", deleteErr)
		}
	}

	return true, nil
}
