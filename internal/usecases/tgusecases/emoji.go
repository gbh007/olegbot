package tgusecases

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *UseCases) EmojiHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}

	botInfo, err := u.repo.GetBot(ctx, u.botID)
	if err != nil {
		return true, fmt.Errorf("emoji handle: bot info: %w", err)
	}

	emoji, ok := botInfo.RandomEmojiWithChance()
	if !ok {
		return false, nil
	}

	_, err = b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
		Reaction: []models.ReactionType{
			{
				Type: models.ReactionTypeTypeEmoji,
				ReactionTypeEmoji: &models.ReactionTypeEmoji{
					Type:  models.ReactionTypeTypeEmoji,
					Emoji: emoji,
				},
			},
		},
	})
	if err != nil {
		return true, fmt.Errorf("emoji handle: set emoji: %w", err)
	}

	return true, nil
}

func (u *UseCases) EmojiCommandHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	if update.Message == nil || update.Message.ReplyToMessage == nil {
		return false, nil
	}

	ok, err := u.commandStrictCheck(ctx, "/emoji", update.Message.Text)
	if err != nil {
		return true, fmt.Errorf("emoji handle: strict check: %w", err)
	}

	if !ok {
		return false, nil
	}

	botInfo, err := u.repo.GetBot(ctx, u.botID)
	if err != nil {
		return true, fmt.Errorf("emoji handle: bot info: %w", err)
	}

	emoji, ok := botInfo.RandomEmoji()
	if !ok {
		return false, nil
	}

	_, err = b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ReplyToMessage.ID,
		Reaction: []models.ReactionType{
			{
				Type: models.ReactionTypeTypeEmoji,
				ReactionTypeEmoji: &models.ReactionTypeEmoji{
					Type:  models.ReactionTypeTypeEmoji,
					Emoji: emoji,
				},
			},
		},
	})
	if err != nil {
		return true, fmt.Errorf("emoji handle: send message: %w", err)
	}

	_, _ = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})

	return true, nil
}
