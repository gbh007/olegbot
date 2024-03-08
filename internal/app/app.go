package app

import (
	"app/internal/controllers/cmscontroller"
	"app/internal/controllers/tgcontroller"
	"app/internal/dataproviders/quote"
	"app/internal/usecases/cmsusecases"
	"app/internal/usecases/tgusecases"
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/vrischmann/envconfig"
)

type appConfig struct {
	Token string
	Bot   struct {
		Name string `envconfig:"optional"`
		Tag  string `envconfig:"optional"`
	} `envconfig:"optional"`
	Repo  string
	Addr  string `envconfig:"optional"`
	Texts struct {
		QuoteAdded  string `envconfig:"optional"`
		QuoteExists string `envconfig:"optional"`
	} `envconfig:"optional"`
	CMS struct {
		Login    string `envconfig:"optional"`
		Password string `envconfig:"optional"`
	} `envconfig:"optional"`
	Debug bool `envconfig:"optional"`
}

type controller interface {
	Serve(ctx context.Context) error
}

type App struct {
	controllers []controller
	logger      *slog.Logger
}

func New(logger *slog.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (a *App) Init(ctx context.Context) error {

	cfg := new(appConfig)

	err := envconfig.Init(cfg)
	if err != nil {
		return fmt.Errorf("app: init: envconfig: %w", err)
	}

	repo := quote.New()

	err = repo.Load(ctx, cfg.Repo)
	if err != nil {
		return fmt.Errorf("app: init: repository: %w", err)
	}

	a.controllers = append(
		a.controllers,
		tgcontroller.New(
			tgcontroller.Config{
				Token:    cfg.Token,
				BotName:  cfg.Bot.Name,
				BotTag:   cfg.Bot.Tag,
				UseCases: tgusecases.New(repo),
				Texts: tgcontroller.Texts{
					QuoteAdded:  cfg.Texts.QuoteAdded,
					QuoteExists: cfg.Texts.QuoteExists,
				},
			},
		),
		cmscontroller.New(
			cmscontroller.Config{
				HTTPAddr: cfg.Addr,
				CMSLogin: cfg.CMS.Login,
				CMSPass:  cfg.CMS.Password,
				Debug:    cfg.Debug,
			},
			cmsusecases.New(repo),
		),
	)

	return nil
}

func (a *App) Serve(ctx context.Context) error {
	wg := new(sync.WaitGroup)
	wg.Add(len(a.controllers))

	for _, c := range a.controllers {
		go func(c controller) {
			defer wg.Done()

			err := c.Serve(ctx)
			if err != nil {
				a.logger.Error("controller serve error", slog.Any("error", err))
			}
		}(c)
	}

	wg.Wait()

	return nil
}
