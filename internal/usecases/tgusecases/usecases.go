package tgusecases

import (
	"app/internal/domain"
	"context"
	"log/slog"
	"strings"
)

type repository interface {
	Quotes(ctx context.Context, botID int64) ([]domain.Quote, error)
	AddQuote(ctx context.Context, botID int64, text string, userID, chatID int64) error
	IsModerator(ctx context.Context, botID int64, userID int64) (bool, error)
	QuoteExists(ctx context.Context, botID int64, text string) (bool, error)
	GetBot(ctx context.Context, botID int64) (domain.Bot, error)
}

type UseCases struct {
	repo  repository
	botID int64

	// FIXME: отрефакторить и убрать отсюда
	logger *slog.Logger
	debug  bool
}

func New(
	repo repository,
	botID int64,
	logger *slog.Logger,
	debug bool,
) *UseCases {
	return &UseCases{
		repo:   repo,
		botID:  botID,
		logger: logger,
		debug:  debug,
	}
}

func (u *UseCases) commandStrictCheck(ctx context.Context, cmd, msg string) (bool, error) {
	bot, err := u.repo.GetBot(ctx, u.botID)
	if err != nil {
		return false, err
	}

	if bot.Tag == "" || !strings.Contains(msg, "@") {
		return strings.HasPrefix(msg, cmd), nil
	}

	ok := strings.HasPrefix(msg, cmd+"@"+bot.Tag)

	return ok, nil
}
