package app

import (
	"app/internal/controller"
	"app/internal/repository"
	"app/internal/usecases"
	"context"
	"fmt"

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
}

type App struct {
	controller *controller.Controller
}

func New() *App {
	return new(App)
}

func (a *App) Init(ctx context.Context) error {

	cfg := new(appConfig)

	err := envconfig.Init(cfg)
	if err != nil {
		return fmt.Errorf("app: init: envconfig: %w", err)
	}

	repo := repository.New()

	err = repo.Load(ctx, cfg.Repo)
	if err != nil {
		return fmt.Errorf("app: init: repository: %w", err)
	}

	a.controller = controller.New(
		controller.Config{
			Token:    cfg.Token,
			BotName:  cfg.Bot.Name,
			BotTag:   cfg.Bot.Tag,
			HTTPAddr: cfg.Addr,
			UseCases: usecases.New(repo),
			Texts: controller.Texts{
				QuoteAdded:  cfg.Texts.QuoteAdded,
				QuoteExists: cfg.Texts.QuoteExists,
			},
		},
	)

	return nil
}

func (a *App) Serve(ctx context.Context) error {
	err := a.controller.Serve(ctx)
	if err != nil {
		return fmt.Errorf("app: serve: %w", err)
	}

	return nil
}
