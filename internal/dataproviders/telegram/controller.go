package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type handler func(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)

type useCases interface {
	EmojiHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	WhoHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	QuoteHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	AddQuoteHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	CommentHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	SelfHandle(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)
	AccessMiddleware() bot.Middleware
}

type Controller struct {
	tgToken  string
	botID    int64
	useCases useCases
	handlers []handler

	logger *slog.Logger
	debug  bool

	cancel func()
	b      *bot.Bot

	wg sync.WaitGroup
}

func New(
	token string,
	botID int64,
	useCases useCases,
	logger *slog.Logger,
	debug bool,
) *Controller {
	c := &Controller{
		tgToken:  token,
		botID:    botID,
		useCases: useCases,
		logger:   logger,
		debug:    debug,
		cancel:   func() {},
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
	ctx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	c.cancel = cancel

	c.wg.Add(1)
	defer c.wg.Done()

	middlewares := make([]bot.Middleware, 0, 3)
	middlewares = append(middlewares, c.counterMiddleware())
	middlewares = append(middlewares, c.useCases.AccessMiddleware())
	middlewares = append(middlewares, c.loggerMiddleware())

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

func (c *Controller) Stop(ctx context.Context) {
	c.cancel()
	c.wg.Wait()
}

func (c *Controller) loggerMiddleware() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			if c.debug {
				data, _ := json.Marshal(update)

				c.logger.DebugContext(
					ctx, "tg bot update",
					slog.Int64("bot_id", c.botID),
					slog.String("data", string(data)),
				)
			}

			next(ctx, bot, update)
		}
	}
}
