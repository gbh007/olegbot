package main

import (
	"app/internal/app"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	}))

	app := app.New(logger)
	logger.Info("app created")

	err := app.Init(ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("app init")

	err = app.Serve(ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("app exit")
}
