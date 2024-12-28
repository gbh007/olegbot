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

func (u *UseCases) randomSticker(ctx context.Context) (string, error) {
	stickers, err := u.repo.Stickers(ctx, u.botID)
	if err != nil {
		return "", fmt.Errorf("use case: random sticker: %w", err)
	}

	if len(stickers) == 0 {
		return "", fmt.Errorf("use case: random sticker: no stickers")
	}

	return stickers[rand.Intn(len(stickers))].FileID, nil
}

func (u *UseCases) addSticker(ctx context.Context, fileID string, userID, chatID int64) error {
	ok, err := u.repo.IsModerator(ctx, u.botID, userID)
	if err != nil {
		return fmt.Errorf("use case: add sticker: %w", err)
	}

	if !ok {
		return fmt.Errorf("use case: add sticker: %w", domain.PermissionDeniedError)
	}

	exists, err := u.repo.StickerExists(ctx, u.botID, fileID)
	if err != nil {
		return fmt.Errorf("use case: add sticker: %w", err)
	}

	if exists {
		return fmt.Errorf("use case: add sticker: %w", domain.AlreadyExistsError)
	}

	err = u.repo.AddSticker(ctx, domain.Sticker{
		BotID:           u.botID,
		FileID:          fileID,
		CreatorID:       userID,
		CreatedInChatID: chatID,
	})
	if err != nil {
		return fmt.Errorf("use case: add sticker: %w", err)
	}

	return nil
}

func (u *UseCases) StickerHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil || update.Message.ReplyToMessage == nil {
		return false, nil
	}

	ok, err := u.commandStrictCheck(ctx, "/sticker", update.Message.Text)
	if err != nil {
		return true, fmt.Errorf("sticker handle: strict check: %w", err)
	}

	if !ok {
		return false, nil
	}

	fileID, err := u.randomSticker(ctx)
	if err != nil {
		return true, fmt.Errorf("sticker handle: %w", err)
	}

	_, err = b.SendSticker(ctx, &bot.SendStickerParams{
		ChatID: update.Message.Chat.ID,
		Sticker: &models.InputFileString{
			Data: fileID,
		},
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ReplyToMessage.ID,
			ChatID:    update.Message.Chat.ID,
		},
	})
	if err != nil {
		return true, fmt.Errorf("sticker handle: send message: %w", err)
	}

	_, _ = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})

	return true, nil
}

func (u *UseCases) AddStickerHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}

	ok, err := u.commandStrictCheck(ctx, "/add_sticker", update.Message.Text)
	if err != nil {
		return true, fmt.Errorf("add sticker handle: strict check: %w", err)
	}

	if !ok {
		return false, nil
	}

	splitIndex := strings.Index(update.Message.Text, " ")
	if splitIndex == -1 && (update.Message.ReplyToMessage == nil || update.Message.ReplyToMessage.Sticker == nil) {
		return true, fmt.Errorf("add sticker invalid syntax")
	}

	var (
		fileID  string
		replyTo int
	)

	switch {
	case update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.Sticker != nil:
		fileID = update.Message.ReplyToMessage.Sticker.FileID
		replyTo = update.Message.ReplyToMessage.ID
	default:
		fileID = update.Message.Text[splitIndex+1:]
		replyTo = update.Message.ID
	}

	err = u.addSticker(ctx, fileID, update.Message.From.ID, update.Message.Chat.ID)
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
			return true, fmt.Errorf("add sticker handle: send message: %w", sendErr)
		}

	case err != nil:
		return true, fmt.Errorf("add sticker handle: %w", err)

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
			return true, fmt.Errorf("add sticker handle: send message: %w", sendErr)
		}
	}

	if update.Message.ReplyToMessage != nil {
		_, deleteErr := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		})
		if deleteErr != nil {
			return true, fmt.Errorf("add sticker handle: delete message: %w", deleteErr)
		}
	}

	return true, nil
}
