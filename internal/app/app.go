package app

import (
	"app/internal/controllers/cmscontroller"
	"app/internal/controllers/tgcontroller"
	"app/internal/dataproviders/postgresql"
	"app/internal/usecases/cmsusecases"
	"app/internal/usecases/tgusecases"
	"context"
	"fmt"
	"log/slog"

	"github.com/vrischmann/envconfig"
)

type appConfig struct {
	Repo string
	Addr string `envconfig:"optional"`
	CMS  struct {
		Login    string `envconfig:"optional"`
		Password string `envconfig:"optional"`
	} `envconfig:"optional"`
	Debug bool `envconfig:"optional"`
}

type App struct {
	logger *slog.Logger

	tgController  *tgcontroller.Controller
	cmsController *cmscontroller.Controller
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

	repo, err := postgresql.New(ctx, cfg.Repo, 10, a.logger, cfg.Debug)
	if err != nil {
		return fmt.Errorf("app: init: repository: %w", err)
	}

	botInfo, err := repo.GetBot(ctx, 1) // FIXME: после рефакторинга не хардкодить
	if err != nil {
		return fmt.Errorf("app: init: get default bot: %w", err)
	}

	a.tgController = tgcontroller.New(
		tgcontroller.Config{
			Token: botInfo.Token,
			UseCases: tgusecases.New(
				repo,
				1, // FIXME: после рефакторинга не хардкодить
			),
		},
	)

	a.cmsController = cmscontroller.New(
		cmscontroller.Config{
			HTTPAddr: cfg.Addr,
			CMSLogin: cfg.CMS.Login,
			CMSPass:  cfg.CMS.Password,
			Debug:    cfg.Debug,
		},
		cmsusecases.New(repo),
		a.tgController, // FIXME: это конечно дич, но пока так проще.
	)

	return nil
}

func (a *App) Serve(ctx context.Context) error {
	go func() {
		err := a.cmsController.Serve(ctx)
		if err != nil {
			a.logger.Error("cms controller serve error", slog.Any("error", err))
		}
	}()

	err := a.tgController.Serve(ctx)
	if err != nil {
		return fmt.Errorf("app: serve: %w", err)
	}

	return nil
}
