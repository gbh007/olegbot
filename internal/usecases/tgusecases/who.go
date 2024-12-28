package tgusecases

import (
	"bytes"
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

	buff := bytes.Buffer{}

	buff.WriteString(">>> message\n")
	buff.WriteString(fmt.Sprintf("user id: %d chat id: %d\n", update.Message.From.ID, update.Message.Chat.ID))

	if update.Message.Animation != nil {
		buff.WriteString(fmt.Sprintf("animation file_id:%s\n", update.Message.Animation.FileID))
	}

	if update.Message.Sticker != nil {
		buff.WriteString(fmt.Sprintf("sticker file_id:%s\n", update.Message.Sticker.FileID))
	}

	if update.Message.ReplyToMessage != nil {
		buff.WriteString(">>> reply\n")
		buff.WriteString(fmt.Sprintf("user id: %d chat id: %d\n", update.Message.ReplyToMessage.From.ID, update.Message.ReplyToMessage.Chat.ID))

		if update.Message.ReplyToMessage.Animation != nil {
			buff.WriteString(fmt.Sprintf("animation file_id:%s\n", update.Message.ReplyToMessage.Animation.FileID))
		}

		if update.Message.ReplyToMessage.Sticker != nil {
			buff.WriteString(fmt.Sprintf("sticker file_id:%s\n", update.Message.ReplyToMessage.Sticker.FileID))
		}
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   buff.String(),
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
