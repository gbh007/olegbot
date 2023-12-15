package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type handler func(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error)

type useCases interface {
	RandomQuote(context.Context) (string, error)
	AddQuote(ctx context.Context, text string, userID, chatID int64) error
}

type Controller struct {
	hasBotName bool
	botName    string

	hasBotTag bool
	botTag    string
	tgToken   string

	useCases useCases

	handlers []handler
}

func New(
	botName string,
	botTag string,
	tgToken string,

	useCases useCases,
) *Controller {
	c := &Controller{
		hasBotName: botName != "",
		botName:    strings.ToLower(botName),

		hasBotTag: botTag != "",
		botTag:    strings.ToLower(botTag),

		tgToken: tgToken,

		useCases: useCases,
	}

	c.handlers = append(
		c.handlers,
		c.commentHandle,
		c.quoteHandle,
		c.addQuoteHandle,
		c.whoHandle,
		c.selfHandle,
	)

	return c
}

func (c *Controller) Serve(ctx context.Context) error {
	opts := []bot.Option{
		bot.WithDefaultHandler(c.handler),
	}

	b, err := bot.New(c.tgToken, opts...)
	if err != nil {
		return fmt.Errorf("serve error: %w", err)
	}

	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{
				Command:     "who",
				Description: "кто я, и где",
			},
			{
				Command:     "quote",
				Description: "великая цитата",
			},
			{
				Command:     "comment",
				Description: "сочный комментарий",
			},
			{
				Command:     "add_quote",
				Description: "добавить цитату",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("serve error: set commands: %w", err)
	}

	b.Start(ctx)

	return nil
}
