package cache

import (
	"app/internal/domain"
	"context"
	"log/slog"
)

type dataSource interface {
	Quotes(ctx context.Context, botID int64) ([]domain.Quote, error)
	AddQuote(ctx context.Context, botID int64, text string, userID, chatID int64) error
	IsModerator(ctx context.Context, botID int64, userID int64) (bool, error)
	QuoteExists(ctx context.Context, botID int64, text string) (bool, error)
	GetBot(ctx context.Context, botID int64) (domain.Bot, error)

	UpdateQuoteText(ctx context.Context, id int64, text string) error
	DeleteQuote(ctx context.Context, id int64) error

	Quote(ctx context.Context, id int64) (domain.Quote, error)

	Moderators(ctx context.Context, botID int64) ([]domain.Moderator, error)
	AddModerator(ctx context.Context, botID, userID int64, description string) error
	DeleteModerator(ctx context.Context, botID, userID int64) error

	CreateBot(ctx context.Context, bot domain.Bot) error
	UpdateBot(ctx context.Context, bot domain.Bot) error
	DeleteBot(ctx context.Context, id int64) error
	GetBots(ctx context.Context) ([]domain.Bot, error)
}

type Cache struct {
	origin dataSource
	logger *slog.Logger

	quotes *wrapper[int64, []domain.Quote]
	bots   *wrapper[int64, domain.Bot]
}

func New(origin dataSource, logger *slog.Logger) *Cache {
	return &Cache{
		origin: origin,
		logger: logger,
		quotes: newWrapper[int64, []domain.Quote]("quotes", logger),
		bots:   newWrapper[int64, domain.Bot]("bots", logger),
	}
}
