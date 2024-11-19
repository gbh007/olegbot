package cmscontroller

import (
	"app/internal/controllers/cmscontroller/internal/binds"
	"app/internal/controllers/cmscontroller/internal/static"
	"app/internal/domain"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type useCases interface {
	Quotes(ctx context.Context) ([]domain.Quote, error)

	Quote(ctx context.Context, id int64) (domain.Quote, error)
	DeleteQuote(ctx context.Context, id int64) error
	UpdateQuoteText(ctx context.Context, id int64, text string) error
	AddQuote(ctx context.Context, text string) error

	Moderators(ctx context.Context) ([]domain.Moderator, error)
	AddModerator(ctx context.Context, userID int64, description string) error
	DeleteModerator(ctx context.Context, userID int64) error

	AddQuotes(ctx context.Context, quotes []string) error
}

type botController interface {
	SendAudio(ctx context.Context, chatID int64, filename string, data io.Reader) error
	SendVideo(ctx context.Context, chatID int64, filename string, data io.Reader) error
}

type Config struct {
	HTTPAddr string

	CMSLogin string
	CMSPass  string

	Debug bool
}

type Controller struct {
	httpAddr string

	cmsLogin string
	cmsPass  string
	debug    bool

	useCases      useCases
	botController botController
}

func New(cfg Config, useCases useCases, botController botController) *Controller {
	return &Controller{
		httpAddr: cfg.HTTPAddr,

		cmsLogin: cfg.CMSLogin,
		cmsPass:  cfg.CMSPass,

		debug: cfg.Debug,

		useCases:      useCases,
		botController: botController,
	}
}

func (c *Controller) Serve(ctx context.Context) error {
	echoRouter := echo.New()

	echoRouter.HideBanner = true
	echoRouter.Debug = c.debug
	echoRouter.Validator = binds.Validator{}

	if c.debug {
		echoRouter.Use(middleware.Logger())
	}

	echoRouter.Use(middleware.Recover())
	echoRouter.Use(c.newBaseAuth())

	echoRouter.StaticFS("/", static.StaticDir)

	echoRouter.GET("/api/quotes", c.quotesHandler())

	echoRouter.GET("/api/quote", c.quoteHandler())
	echoRouter.DELETE("/api/quote", c.deleteQuoteHandler())
	echoRouter.POST("/api/quote", c.updateQuoteHandler())
	echoRouter.PUT("/api/quote", c.addQuoteHandler())

	echoRouter.GET("/api/moderators", c.moderatorsHandler())
	echoRouter.DELETE("/api/moderator", c.deleteModeratorHandler())
	echoRouter.PUT("/api/moderator", c.addModeratorHandler())

	echoRouter.POST("/api/ff/quotes", c.ffQuoteHandler())
	echoRouter.POST("/api/ff/media", c.ffMediaHandler())

	echoRouter.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	go func() { // Стоит переписать, пока временная затычка
		<-ctx.Done()
		_ = echoRouter.Shutdown(context.Background())
	}()

	err := echoRouter.Start(c.httpAddr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http serve: %w", err)
	}

	return nil
}
