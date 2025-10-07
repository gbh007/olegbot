package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"app/internal/controllers/cmscontroller"
	"app/internal/controllers/tgcontroller"
	"app/internal/dataproviders/cache"
	"app/internal/dataproviders/deepseek"
	"app/internal/dataproviders/ollama"
	"app/internal/dataproviders/openai"
	"app/internal/dataproviders/postgresql"
	"app/internal/usecases/cmsusecases"

	"github.com/BurntSushi/toml"
)

type appConfig struct {
	Repo string `toml:"repo"`
	Addr string `toml:"addr"`
	Llm  struct {
		Addr    string        `toml:"addr"`
		Model   string        `toml:"model"`
		Token   string        `toml:"token"`
		Type    string        `toml:"type"`
		Timeout time.Duration `toml:"timeout"`
	} `toml:"llm"`
	CMS struct {
		StaticDirPath string `toml:"static_dir_path"`
		Login         string `toml:"login"`
		Password      string `toml:"password"`
	} `toml:"cms"`
	Debug bool `toml:"debug"`
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
	cfg := appConfig{}

	_, err := toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		return fmt.Errorf("app: init: config: %w", err)
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
	case cfg.Llm.Type == "openai" && cfg.Llm.Token != "" && cfg.Llm.Addr != "" && cfg.Llm.Model != "":
		llmProvider, err = openai.New(ctx, a.logger, cfg.Llm.Addr, cfg.Llm.Token, cfg.Llm.Model)
		if err != nil {
			return fmt.Errorf("app: init: openai: %w", err)
		}
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
		cfg.Llm.Timeout,
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
