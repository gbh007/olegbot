package tgcontroller

import (
	"context"
	"io"
)

func (c *Controller) SendAudio(ctx context.Context, botID, chatID int64, filename string, data io.Reader) error {
	c.botsMutex.RLock()
	defer c.botsMutex.RUnlock()

	bot, ok := c.bots[botID]
	if !ok {
		return botNotRunningErr
	}

	return bot.SendAudio(ctx, chatID, filename, data)
}

func (c *Controller) SendVideo(ctx context.Context, botID, chatID int64, filename string, data io.Reader) error {
	c.botsMutex.RLock()
	defer c.botsMutex.RUnlock()

	bot, ok := c.bots[botID]
	if !ok {
		return botNotRunningErr
	}

	return bot.SendVideo(ctx, chatID, filename, data)
}

func (c *Controller) SendImage(ctx context.Context, botID, chatID int64, filename string, data io.Reader) error {
	c.botsMutex.RLock()
	defer c.botsMutex.RUnlock()

	bot, ok := c.bots[botID]
	if !ok {
		return botNotRunningErr
	}

	return bot.SendImage(ctx, chatID, filename, data)
}

func (c *Controller) SendText(ctx context.Context, botID, chatID int64, text string) error {
	c.botsMutex.RLock()
	defer c.botsMutex.RUnlock()

	bot, ok := c.bots[botID]
	if !ok {
		return botNotRunningErr
	}

	return bot.SendText(ctx, chatID, text)
}
