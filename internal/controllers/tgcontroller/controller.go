package tgcontroller

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Config struct {
	Token string

	BotName      string
	BotTag       string
	Tags         []string
	AllowedChats []int64

	UseCases useCases

	Texts Texts
}

type Texts struct {
	QuoteAdded  string
	QuoteExists string
}

type handler func(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)

type useCases interface {
	RandomQuote(context.Context) (string, error)
	AddQuote(ctx context.Context, text string, userID, chatID int64) error
}

type Controller struct {
	tags         []string
	allowedChats []int64

	tgToken string

	useCases useCases

	handlers []handler

	texts Texts
}

func New(cfg Config) *Controller {
	tags := make([]string, 0, len(cfg.Tags)+2)

	if cfg.BotName != "" {
		tags = append(tags, strings.ToLower(cfg.BotName))
	}

	if cfg.BotTag != "" {
		tags = append(tags, strings.ToLower(cfg.BotTag))
	}

	for _, tag := range cfg.Tags {
		if tag != "" {
			tags = append(tags, strings.ToLower(tag))
		}
	}

	c := &Controller{
		tags:         tags,
		allowedChats: cfg.AllowedChats,

		tgToken: cfg.Token,

		useCases: cfg.UseCases,

		texts: cfg.Texts,
	}

	if c.texts.QuoteAdded == "" {
		c.texts.QuoteAdded = "✅ quote added"
	}

	if c.texts.QuoteExists == "" {
		c.texts.QuoteExists = "❌ quote already exists"
	}

	c.handlers = append(
		c.handlers,
		c.handleWrapper(c.commentHandle, "comment"),
		c.handleWrapper(c.quoteHandle, "quote"),
		c.handleWrapper(c.addQuoteHandle, "add_quote"),
		c.handleWrapper(c.whoHandle, "who"),
		c.handleWrapper(c.selfHandle, "self"),
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

	b.Start(ctx)

	return nil
}
