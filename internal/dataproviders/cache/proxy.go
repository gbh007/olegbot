package cache

import (
	"app/internal/domain"
	"context"
)

func (c *Cache) Quotes(ctx context.Context, botID int64) ([]domain.Quote, error) {
	quotes, ok := c.quotes.Get(botID)
	if ok {
		return quotes, nil
	}

	quotes, err := c.origin.Quotes(ctx, botID)
	if err == nil {
		c.quotes.Set(botID, quotes)
	}

	return quotes, err
}

func (c *Cache) AddQuote(ctx context.Context, botID int64, text string, userID, chatID int64) error {
	c.quotes.Del(botID)

	return c.origin.AddQuote(ctx, botID, text, userID, chatID)
}

func (c *Cache) IsModerator(ctx context.Context, botID int64, userID int64) (bool, error) {
	return c.origin.IsModerator(ctx, botID, userID)
}

func (c *Cache) QuoteExists(ctx context.Context, botID int64, text string) (bool, error) {
	return c.origin.QuoteExists(ctx, botID, text)
}

func (c *Cache) GetBot(ctx context.Context, botID int64) (domain.Bot, error) {
	bot, ok := c.bots.Get(botID)
	if ok {
		return bot, nil
	}

	bot, err := c.origin.GetBot(ctx, botID)
	if err == nil {
		c.bots.Set(botID, bot)
	}

	return bot, err
}

func (c *Cache) UpdateQuoteText(ctx context.Context, id int64, text string) error {
	c.quotes.Clear()

	return c.origin.UpdateQuoteText(ctx, id, text)
}

func (c *Cache) DeleteQuote(ctx context.Context, id int64) error {
	c.quotes.Clear()

	return c.origin.DeleteQuote(ctx, id)
}

func (c *Cache) Quote(ctx context.Context, id int64) (domain.Quote, error) {
	return c.origin.Quote(ctx, id)
}

func (c *Cache) Moderators(ctx context.Context, botID int64) ([]domain.Moderator, error) {
	return c.origin.Moderators(ctx, botID)
}

func (c *Cache) AddModerator(ctx context.Context, botID, userID int64, description string) error {
	return c.origin.AddModerator(ctx, botID, userID, description)
}

func (c *Cache) DeleteModerator(ctx context.Context, botID, userID int64) error {
	return c.origin.DeleteModerator(ctx, botID, userID)
}

func (c *Cache) CreateBot(ctx context.Context, bot domain.Bot) error {
	return c.origin.CreateBot(ctx, bot)
}

func (c *Cache) UpdateBot(ctx context.Context, bot domain.Bot) error {
	c.bots.Del(bot.ID)

	return c.origin.UpdateBot(ctx, bot)
}

func (c *Cache) DeleteBot(ctx context.Context, id int64) error {
	c.bots.Del(id)

	return c.origin.DeleteBot(ctx, id)
}

func (c *Cache) GetBots(ctx context.Context) ([]domain.Bot, error) {
	return c.origin.GetBots(ctx)
}
