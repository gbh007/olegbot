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
		Name string
		Tag  string
	}
	Repo string
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
		cfg.Bot.Name, cfg.Bot.Tag, cfg.Token,
		usecases.New(repo),
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
