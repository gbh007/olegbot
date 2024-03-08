package cmscontroller

import (
	"app/internal/controllers/cmscontroller/internal/static"
	"app/internal/domain"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type useCases interface {
	Quotes(ctx context.Context) ([]domain.Quote, error)
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

	useCases useCases
}

func New(cfg Config, useCases useCases) *Controller {
	return &Controller{
		httpAddr: cfg.HTTPAddr,

		cmsLogin: cfg.CMSLogin,
		cmsPass:  cfg.CMSPass,

		debug: cfg.Debug,

		useCases: useCases,
	}
}

func (c *Controller) Serve(ctx context.Context) error {
	echoRouter := echo.New()

	echoRouter.HideBanner = true
	echoRouter.Debug = c.debug

	if c.debug {
		echoRouter.Use(middleware.Logger())
	}

	echoRouter.Use(middleware.Recover())
	echoRouter.Use(c.newBaseAuth())

	echoRouter.StaticFS("/", static.StaticDir)

	echoRouter.GET("/api/quotes", c.quotesHandler())

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
