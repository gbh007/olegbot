package telegram

import (
	"context"
	"io"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Controller) SendAudio(ctx context.Context, chatID int64, filename string, data io.Reader) error {
	_, err := c.b.SendAudio(ctx, &bot.SendAudioParams{
		ChatID: chatID,
		Audio: &models.InputFileUpload{
			Filename: filename,
			Data:     data,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) SendVideo(ctx context.Context, chatID int64, filename string, data io.Reader) error {
	_, err := c.b.SendVideo(ctx, &bot.SendVideoParams{
		ChatID: chatID,
		Video: &models.InputFileUpload{
			Filename: filename,
			Data:     data,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) SendImage(ctx context.Context, chatID int64, filename string, data io.Reader) error {
	_, err := c.b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: chatID,
		Photo: &models.InputFileUpload{
			Filename: filename,
			Data:     data,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) SendText(ctx context.Context, chatID int64, text string) error {
	_, err := c.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return err
	}

	return nil
}
