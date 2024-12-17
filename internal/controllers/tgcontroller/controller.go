package tgcontroller

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Config struct {
	Token string

	BotName      string
	BotTag       string
	AllowedChats []int64

	UseCases useCases
}

type handler func(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)

type useCases interface {
	EmojiHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	WhoHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	QuoteHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	AddQuoteHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	CommentHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	SelfHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
}

type Controller struct {
	allowedChats []int64

	tgToken string

	useCases useCases

	handlers []handler

	b *bot.Bot
}

func New(cfg Config) *Controller {
	c := &Controller{
		allowedChats: cfg.AllowedChats,

		tgToken: cfg.Token,

		useCases: cfg.UseCases,
	}

	c.handlers = append(
		c.handlers,
		c.handleWrapper(c.useCases.CommentHandle, "comment"),
		c.handleWrapper(c.useCases.QuoteHandle, "quote"),
		c.handleWrapper(c.useCases.AddQuoteHandle, "add_quote"),
		c.handleWrapper(c.useCases.WhoHandle, "who"),
		c.handleWrapper(c.useCases.SelfHandle, "self"),
		c.handleWrapper(c.useCases.EmojiHandle, "emoji"),
	)

	return c
}

func (c *Controller) Serve(ctx context.Context) error {
	middlewares := make([]bot.Middleware, 0, 2)
	middlewares = append(middlewares, c.counterMiddleware())

	if len(c.allowedChats) > 0 {
		middlewares = append(middlewares, c.accessMiddleware())
	}

	opts := []bot.Option{
		bot.WithMiddlewares(middlewares...),
		bot.WithDefaultHandler(c.handler),
	}

	b, err := bot.New(c.tgToken, opts...)
	if err != nil {
		return fmt.Errorf("serve error: %w", err)
	}

	err = c.setBotCommands(ctx, b)
	if err != nil {
		return fmt.Errorf("serve error: %w", err)
	}

	c.b = b

	b.Start(ctx)

	return nil
}
