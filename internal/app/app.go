package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"app/internal/controllers/cmscontroller"
	"app/internal/controllers/tgcontroller"
	"app/internal/dataproviders/cache"
	"app/internal/dataproviders/deepseek"
	"app/internal/dataproviders/ollama"
	"app/internal/dataproviders/postgresql"
	"app/internal/usecases/cmsusecases"

	"github.com/vrischmann/envconfig"
)

type appConfig struct {
	Repo string
	Addr string `envconfig:"optional"`
	Llm  struct {
		Addr  string `envconfig:"optional"`
		Model string `envconfig:"optional"`
		Token string `envconfig:"optional"`
		Type  string `envconfig:"optional"`
	} `envconfig:"optional"`
	CMS struct {
		StaticDirPath string `envconfig:"optional,"`
		Login         string `envconfig:"optional"`
		Password      string `envconfig:"optional"`
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

	logLevel := slog.LevelInfo
	if cfg.Debug {
		logLevel = slog.LevelDebug
	}

	a.logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: cfg.Debug,
		Level:     logLevel,
	}))

	repo, err := postgresql.New(ctx, cfg.Repo, 10, a.logger, cfg.Debug)
	if err != nil {
		return fmt.Errorf("app: init: repository: %w", err)
	}

	var llmProvider tgcontroller.Llm

	switch {
	case cfg.Llm.Type == "deepseek" && cfg.Llm.Token != "":
		llmProvider, err = deepseek.New(ctx, a.logger, cfg.Llm.Token)
		if err != nil {
			return fmt.Errorf("app: init: deepseek: %w", err)
		}
	case cfg.Llm.Type == "ollama" && cfg.Llm.Addr != "" && cfg.Llm.Model != "":
		llmProvider, err = ollama.New(ctx, a.logger, cfg.Llm.Addr, cfg.Llm.Model)
		if err != nil {
			return fmt.Errorf("app: init: ollama: %w", err)
		}
	}

	cachedRepo := cache.New(repo, a.logger)

	a.tgController = tgcontroller.New(
		cachedRepo,
		llmProvider,
		a.logger,
		cfg.Debug,
	)

	a.cmsController = cmscontroller.New(
		cmscontroller.Config{
			HTTPAddr:      cfg.Addr,
			CMSLogin:      cfg.CMS.Login,
			CMSPass:       cfg.CMS.Password,
			Debug:         cfg.Debug,
			StaticDirPath: cfg.CMS.StaticDirPath,
		},
		cmsusecases.New(cachedRepo),
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
