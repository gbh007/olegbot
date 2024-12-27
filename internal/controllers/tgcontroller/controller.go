package tgcontroller

import (
	"app/internal/dataproviders/telegram"
	"app/internal/domain"
	"app/internal/usecases/tgusecases"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

var (
	botNotRunningErr = errors.New("bot is not running")
	botNotEnabledErr = errors.New("bot is not enabled")
)

type Controller struct {
	// FIXME: это очень плохо, надо отрефакторить (включая юзкейсы), чтобы перейти на интерфейсы
	bots      map[int64]*telegram.Controller
	botsMutex sync.RWMutex

	logger *slog.Logger

	repo repo
}

type repo interface {
	GetBots(ctx context.Context) ([]domain.Bot, error)

	Quotes(ctx context.Context, botID int64) ([]domain.Quote, error)
	AddQuote(ctx context.Context, botID int64, text string, userID, chatID int64) error
	IsModerator(ctx context.Context, botID int64, userID int64) (bool, error)
	QuoteExists(ctx context.Context, botID int64, text string) (bool, error)
	GetBot(ctx context.Context, botID int64) (domain.Bot, error)
}

func New(repo repo, logger *slog.Logger) *Controller {
	c := &Controller{
		bots:   make(map[int64]*telegram.Controller),
		repo:   repo,
		logger: logger,
	}

	return c
}

func (c *Controller) Serve(ctx context.Context) error {
	bots, err := c.repo.GetBots(ctx)
	if err != nil {
		return fmt.Errorf("get bots: %w", err)
	}

	for _, bot := range bots {
		if bot.Enabled {
			c.startBot(ctx, bot)
		}
	}

	<-ctx.Done()

	c.botsMutex.Lock()

	for _, bot := range c.bots {
		bot.Stop(context.TODO())
	}

	c.botsMutex.Unlock()

	return nil
}

func (c *Controller) startBot(ctx context.Context, bot domain.Bot) {
	// FIXME: это очень плохо, надо отрефакторить (включая юзкейсы), чтобы перейти на интерфейсы
	bc := telegram.New(bot.Token, tgusecases.New(c.repo, bot.ID))

	c.botsMutex.Lock()
	_, exists := c.bots[bot.ID]
	if !exists {
		c.bots[bot.ID] = bc
	}
	c.botsMutex.Unlock()

	if exists {
		return
	}

	go func() {
		err := bc.Serve(ctx)
		if err != nil {
			c.logger.ErrorContext(ctx, err.Error())
		}

		c.botsMutex.Lock()
		delete(c.bots, bot.ID)
		c.botsMutex.Unlock()
	}()
}

func (c *Controller) StartBot(ctx context.Context, botID int64) error {
	bot, err := c.repo.GetBot(ctx, botID)
	if err != nil {
		return err
	}

	if bot.Enabled {
		return botNotEnabledErr
	}

	c.startBot(ctx, bot)

	return nil
}

func (c *Controller) StopBot(ctx context.Context, botID int64) error {
	c.botsMutex.Lock()

	bot, ok := c.bots[botID]
	if ok {
		bot.Stop(ctx)
	}

	c.botsMutex.Unlock()

	if !ok {
		return botNotRunningErr
	}

	return nil
}

func (c *Controller) RunningBots(ctx context.Context) ([]int64, error) {
	c.botsMutex.RLock()

	ids := make([]int64, len(c.bots))
	for id := range c.bots {
		ids = append(ids, id)
	}

	c.botsMutex.RUnlock()

	return ids, nil
}
