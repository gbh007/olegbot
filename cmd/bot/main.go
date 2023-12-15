package main

import (
	"app/internal/app"
	"context"
	"log"
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

	app := app.New()
	log.Println("app created")

	err := app.Init(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("app init")

	err = app.Serve(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("app exit")
}
