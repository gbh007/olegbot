package app

import (
	"app/internal/controllers/cmscontroller"
	"app/internal/controllers/tgcontroller"
	"app/internal/dataproviders/postgresql"
	"app/internal/domain"
	"app/internal/usecases/cmsusecases"
	"app/internal/usecases/tgusecases"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/vrischmann/envconfig"
)

type appConfig struct {
	Token string
	Bot   struct {
		Name         string   `envconfig:"optional"`
		Tag          string   `envconfig:"optional"`
		Tags         []string `envconfig:"optional"`
		AllowedChats []int64  `envconfig:"optional"`
	} `envconfig:"optional"`
	Repo string
	Addr string `envconfig:"optional"`
	CMS  struct {
		Login    string `envconfig:"optional"`
		Password string `envconfig:"optional"`
	} `envconfig:"optional"`
	Debug bool `envconfig:"optional"`
	Emoji struct {
		List   []string `envconfig:"optional"`
		Chance float32  `envconfig:"optional"`
	} `envconfig:"optional"`
}

func (cfg appConfig) toBot() domain.Bot {
	tags := make([]string, 0, len(cfg.Bot.Tags)+2)

	if cfg.Bot.Name != "" {
		tags = append(tags, strings.ToLower(cfg.Bot.Name))
	}

	if cfg.Bot.Tag != "" {
		tags = append(tags, strings.ToLower(cfg.Bot.Tag))
	}

	for _, tag := range cfg.Bot.Tags {
		if tag != "" {
			tags = append(tags, strings.ToLower(tag))
		}
	}

	return domain.Bot{
		EmojiList:    cfg.Emoji.List,
		EmojiChance:  cfg.Emoji.Chance,
		Tags:         tags,
		Name:         cfg.Bot.Name,
		Tag:          cfg.Bot.Tag,
		Token:        cfg.Token,
		AllowedChats: cfg.Bot.AllowedChats,
	}
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

	repo := postgresql.New(cfg.toBot())

	err = repo.Load(ctx, cfg.Repo)
	if err != nil {
		return fmt.Errorf("app: init: repository: %w", err)
	}

	fmt.Println(cfg.Bot.AllowedChats)

	a.tgController = tgcontroller.New(
		tgcontroller.Config{
			Token:        cfg.Token,
			BotName:      cfg.Bot.Name,
			BotTag:       cfg.Bot.Tag,
			AllowedChats: cfg.Bot.AllowedChats,
			UseCases: tgusecases.New(
				repo,
				cfg.Emoji.List,
				cfg.Emoji.Chance,
				cfg.Bot.Tags,
				cfg.Bot.Name, cfg.Bot.Tag,
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
